package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qa-dev/jsonwire-grid/jsonwire"
)

var (
	hubUrl      *string
	level       *int
	durationStr *string
)

// todo: prototype
func main() {
	hubUrl = flag.String("hub", "http://127.0.0.1:4444", "address of hub, default http://127.0.0.1:4444")
	level = flag.Int("level", 100, "count parallell conections")
	durationStr = flag.String("duration", "60s", "duration of test, string format ex. 12m, see time.ParseDuration()")
	mockMaxDuration := flag.Int("mockMaxDuration", 500, "request duration [0 <=duration], default 0")
	mockStartPort := flag.Int("mockStartPort", 5000, "mockStartPort")
	flag.Parse()

	duration, err := time.ParseDuration(*durationStr)
	if err != nil {
		log.Fatal("Invalid duration")
	}
	var counter uint64 = 0

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	wg := sync.WaitGroup{}
	isAlive := true
	errChan := make(chan error, *level)

	cmd := exec.Command(
		"webdriver-mock-creator",
		fmt.Sprintf("-hub=%v", *hubUrl),
		fmt.Sprintf("-startPort=%v", *mockStartPort),
		fmt.Sprintf("-maxDuration=%v", *mockMaxDuration),
		fmt.Sprintf("-countNodes=%v", *level),
	)

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 500)

	go func() {
		for i := 1; i <= *level && isAlive; i++ {
			time.Sleep(time.Millisecond * 100)
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					if !isAlive {
						break
					}
					err := runScenario()
					if err != nil {
						errChan <- errors.New("Run scenario, " + err.Error())
					}
					atomic.AddUint64(&counter, 1)
				}
			}()
		}
	}()

	select {
	case <-time.After(duration):
	case err = <-errChan:
	case <-stop:
	}

	isAlive = false

	//wait end all running scenarios
	wg.Wait()

	//wait interrupt child process
	_ = cmd.Process.Signal(os.Interrupt)
	_ = cmd.Wait()

	if err != nil {
		log.Fatalf("Tests failed: %v, ", err)
	}
	log.Printf("Test ok, %v cycles passed", counter)
}

func runScenario() error {
	sessionID, err := createSession()
	if err != nil {
		err = errors.New("Create session, " + err.Error())
		return err
	}
	err = sendAnyRequest(sessionID)
	if err != nil {
		err = errors.New("Send any request, " + err.Error())
		return err
	}
	err = closeSession(sessionID)
	if err != nil {
		err = errors.New("Close session, " + err.Error())
		return err
	}
	return nil
}

func createSession() (sessionID string, err error) {
	resp, err := http.Post(*hubUrl+"/wd/hub/session", "application/json", bytes.NewBuffer([]byte(`{"desiredCapabilities":{"browserName": "firefox"}}`)))
	if err != nil {
		err = errors.New("Send request, " + err.Error())
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("Read response body, " + err.Error())
		return
	}

	var message jsonwire.NewSession
	err = json.Unmarshal(b, &message)
	if err != nil {
		err = errors.New("Unmarshal json, " + err.Error() + ", given response body=[" + string(b) + "]")
		return
	}
	switch {
	case message.SessionID != "":
		sessionID = message.SessionID
	case message.Value.SessionID != "":
		sessionID = message.Value.SessionID
	default:
		err = errors.New("Field`s SessionID is empty")
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("Expected status code 200, actual: " + strconv.Itoa(resp.StatusCode) + ", given response body=[" + string(b) + "]")
		return
	}
	return sessionID, nil
}

func sendAnyRequest(sessionID string) (err error) {
	resp, err := http.Get(*hubUrl + "/wd/hub/session/" + sessionID + "/url")
	if err != nil {
		err = errors.New("Send request, " + err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		err = errors.New("Expected status code 200, actual: " + strconv.Itoa(resp.StatusCode) + ", given response body=[" + string(b) + "]")
		return
	}
	return
}

func closeSession(sessionID string) (err error) {
	req, err := http.NewRequest(http.MethodDelete, *hubUrl+"/wd/hub/session/"+sessionID, nil)
	if err != nil {
		err = errors.New("Create request, " + err.Error())
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.New("Send request, " + err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		err = errors.New("Expected status code 200, actual: " + strconv.Itoa(resp.StatusCode) + ", given response body=[" + string(b) + "]")
		return
	}
	return
}
