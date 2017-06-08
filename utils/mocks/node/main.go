package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	rand.Seed(time.Now().UTC().UnixNano())
	portFlag := flag.Int("port", rand.Intn(1000)+5000, "port default, rand")
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
			<-time.Tick(time.Second)
			err := sendApiProxy()
			if err != nil {
				log.Errorf("Error send [api/proxy], ", err)
			}
		}
	}()

	http.HandleFunc("/wd/hub/session", createSession)
	http.HandleFunc("/wd/hub/session/", useSession)
	http.HandleFunc("/wd/hub/sessions", getSessions)

	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Errorf("Listen serve error, %s", err)
	}
}

// getIpv4 get ip of this server
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

// register register on hub
func register() {
	log.Info("Try register")
	register := jsonwire.Register{
		Configuration:    &jsonwire.Configuration{Host: host, Port: port},
		CapabilitiesList: []jsonwire.Capabilities{{"browserName": "firefox"}},
	}
	b, err := json.Marshal(register)
	if err != nil {
		log.Errorf("Can't encode register json, %s", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, hubUrl+"/grid/register", bytes.NewBuffer(b))
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

// sendApiProxy check "is server know me" and register if server return false
func sendApiProxy() error {
	b := strings.NewReader("{}")
	req, err := http.NewRequest(http.MethodPost, hubUrl+"/grid/api/proxy?id=http://"+host+":"+strconv.Itoa(port), b)
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
