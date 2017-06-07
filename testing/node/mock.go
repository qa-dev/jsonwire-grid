package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os/exec"
	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"encoding/json"
	"regexp"
	"strings"
	"flag"
	"time"
	"net"
	"strconv"
	"bytes"
	"io/ioutil"
	"fmt"
	"errors"
	"math/rand"
)

var currentSessionID, host, hubUrl string
var maxDuration, port int

func main() {
	var err error
	host, err = getIpv4()
	if err != nil {
		log.Fatalf("Can't get host, %s", err)
	}

	hubUrlFlag := flag.String("hub", "http://127.0.0.1:4444", "address of hub, default http://127.0.0.1:4444")
	rand.Seed( time.Now().UTC().UnixNano())
	portFlag := flag.Int("port", rand.Intn(1000) + 5000, "port default, rand")
	maxDurationFlag := flag.Int("maxDuration", 0, "request duration [0 <=duration], default 0")
	flag.Parse()
	hubUrl = *hubUrlFlag
	port = *portFlag
	maxDuration = *maxDurationFlag
	log.Infof("hub url: %v", hubUrl)
	log.Infof("port: %v", port)
	log.Infof("maxDuration: %v", maxDuration)
	register()

	go func() {
		for {
			<- time.Tick(time.Second)
			err := sendApiProxy()
			if err != nil {
				log.Errorf("Error send [api/proxy], ", err)
			}
		}
	}()

	http.HandleFunc("/wd/hub/session", createSession)
	http.HandleFunc("/wd/hub/session/", useSession)
	http.HandleFunc("/wd/hub/sessions", getSessions)

	err = http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		log.Errorf("Listen serve error, %s", err)
	}
}

// getIpv4 метод получает ip адрес, на котором стартует данный сервис.
func getIpv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}

func createSession(rw http.ResponseWriter, r *http.Request) {
	if maxDuration > 0 {
		rand.Seed( time.Now().UTC().UnixNano())
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(maxDuration)))
	}
	rw.Header().Set("Accept", "application/json")
	rw.Header().Set("Accept-charset", "utf-8")

	if r.Method != http.MethodPost {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if currentSessionID != "" {
		errorMassage := "Session already exists"
		log.Error(errorMassage)
		rw.WriteHeader(http.StatusInternalServerError)
		responseMessage := &jsonwire.Message{}
		responseMessage.Status = int(jsonwire.RESPONSE_STATUS_UNKNOWN_ERR)
		responseMessage.Value = errorMassage
		json.NewEncoder(rw).Encode(responseMessage)
		return
	}

	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatalf("Can't generate uuid, %s", err)
	}
	currentSessionID = string(out[:len(out) - 1]) // cut end of line char
	json.NewEncoder(rw).Encode(&jsonwire.Message{SessionId: currentSessionID})
}

func useSession(rw http.ResponseWriter, r *http.Request) {
	if maxDuration > 0 {
		rand.Seed( time.Now().UTC().UnixNano())
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(maxDuration)))
	}
	rw.Header().Set("Accept", "application/json")
	rw.Header().Set("Accept-charset", "utf-8")

	re := regexp.MustCompile(".*/session/([^/]+)(?:/([^/]+))?")
	parsedUrl := re.FindStringSubmatch(r.URL.Path)
	if len(parsedUrl) != 3 {
		errorMessage := "url [" + r.URL.Path + "] parsing error"
		log.Infof(errorMessage)
		http.Error(rw, errorMessage, http.StatusBadRequest)
		return
	}
	sessionId := re.FindStringSubmatch(r.URL.Path)[1]
	responseMessage := &jsonwire.Message{SessionId: sessionId}
	if sessionId != currentSessionID {
		errorMassage := fmt.Sprintf("sessionID '%s' not found", sessionId)
		log.Error(errorMassage)
		rw.WriteHeader(http.StatusNotFound)
		responseMessage.Status = int(jsonwire.RESPONSE_STATUS_UNKNOWN_ERR)
		responseMessage.Value = errorMassage
		json.NewEncoder(rw).Encode(responseMessage)
		return
	}
	if parsedUrl[2] == "" && r.Method == http.MethodDelete {
		currentSessionID = ""
	}
	json.NewEncoder(rw).Encode(responseMessage)
}

func register() {
	log.Info("Try register")
	register := jsonwire.Register{
		Configuration: &jsonwire.Configuration{Host:host, Port: port},
		CapabilitiesList: []jsonwire.Capabilities{{"browserName": "firefox"}},
	}
	b, err := json.Marshal(register)
	if err != nil {
		log.Errorf("Can't encode register json, %s", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, hubUrl + "/grid/register", bytes.NewBuffer(b))
	if err != nil {
		log.Errorf("Can't register, create request error, %s", err)
		return
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-charset", "utf-8")
	err = req.Body.Close()
	if err != nil {
		log.Errorf("Can't register, close request body error, %s", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Can't register, send request error, %s", err)
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Can't register, read response body error, %s", err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		log.Errorf("Can't register, close response body error, %s", err)
		return
	}
	if string(respBytes) != "ok" {
		log.Errorf("Can't register, response body unexpected, %s", respBytes)
		return
	}
	log.Info("Success register")

}

func sendApiProxy() error {
	b := strings.NewReader("{}")
	req, err := http.NewRequest(http.MethodPost, hubUrl + "/grid/api/proxy?id=http://" + host + ":" + strconv.Itoa(port) , b)
	if err != nil {
		err = errors.New(fmt.Sprintf("create request error, %s", err))
		return err
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-charset", "utf-8")
	err = req.Body.Close()
	if err != nil {
		err = errors.New(fmt.Sprintf("close request body error, %s", err))
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("send request error, %s", err))
		return err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("read response body error, %s", err))
		return err
	}
	var respStruct jsonwire.ApiProxy
	err = json.Unmarshal(respBytes, &respStruct)
	if err != nil {
		err = errors.New(fmt.Sprintf("decode json, %s", err))
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		err = errors.New(fmt.Sprintf("close response body error, %s", err))
		return err
	}
	if !respStruct.Success {
		log.Info("Node not registered on hub")
		register()
	}
	return nil
}

func getSessions(rw http.ResponseWriter, r *http.Request) {
	sessions := &jsonwire.Sessions{}
	if currentSessionID != "" {
		sessions.Value = []struct {
			Id           string `json:"id"`
			Capabilities json.RawMessage `json:"capabilities"`
		}{
			{Id: currentSessionID, Capabilities: nil},
		}
	}


	err := json.NewEncoder(rw).Encode(sessions)
	if err != nil {
		err = errors.New("Get sessions error, " + err.Error())
		log.Error(err)
		json.NewEncoder(rw).Encode(&jsonwire.Message{Value: err.Error(), Status: int(jsonwire.RESPONSE_STATUS_UNKNOWN_ERR)})
	}



}




