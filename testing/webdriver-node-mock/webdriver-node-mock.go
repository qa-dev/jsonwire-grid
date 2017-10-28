package main

import (
	"bytes"
	"encoding/json"
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

var currentSessionID, host, hubUrl, id string
var maxDuration, port int

func main() {
	var err error
	host, err = getIpv4()
	if err != nil {
		log.Fatalf("Can't get host, %s", err)
	}

	hubUrlFlag := flag.String("hub", "http://127.0.0.1:4444", "address of hub, default http://127.0.0.1:4444")
	rand.Seed(time.Now().UTC().UnixNano())
	portFlag := flag.Int("port", 5555, "port default, rand")
	maxDurationFlag := flag.Int("maxDuration", 0, "request duration [0 <=duration], default 0")
	flag.Parse()
	hubUrl = *hubUrlFlag
	port = *portFlag
	maxDuration = *maxDurationFlag
	log.Infof("hub url: %v", hubUrl)
	log.Infof("port: %v", port)
	log.Infof("maxDuration: %v", maxDuration)

	id = "http://" + host + ":" + strconv.Itoa(port)
	register()

	go func() {
		for {
			time.Sleep(time.Second)
			err := sendApiProxy()
			if err != nil {
				log.Errorf("Error send [api/proxy], %s", err)
			}
		}
	}()

	http.HandleFunc("/wd/hub/session", createSession)
	http.HandleFunc("/wd/hub/session/", useSession)
	http.HandleFunc("/wd/hub/sessions", getSessions)
	http.HandleFunc("/wd/hub/status", status)

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
		Configuration: &jsonwire.Configuration{
			ID:               id,
			Host:             host,
			Port:             port,
			CapabilitiesList: []jsonwire.Capabilities{{"browserName": "firefox"}}},
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
	b := strings.NewReader(`{id:"` + id + `"}`)
	req, err := http.NewRequest(http.MethodPost, hubUrl+"/grid/api/proxy?id="+id, b)
	if err != nil {
		return fmt.Errorf("create request error, %s", err)
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-charset", "utf-8")
	err = req.Body.Close()
	if err != nil {
		return fmt.Errorf("close request body error, %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request error, %s", err)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body error, %s", err)
	}
	var respStruct jsonwire.APIProxy
	err = json.Unmarshal(respBytes, &respStruct)
	if err != nil {
		return fmt.Errorf("decode json, %s", err)
	}
	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("close response body error, %s", err)
	}
	if !respStruct.Success {
		log.Info("Node not registered on hub")
		log.Info(string(respBytes))
		register()
	}
	return nil
}
