package selenium

import (
	"encoding/json"
	"fmt"
	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	address string
}

type request struct {
	httpRequest *http.Request
}

type HTTPMethod string

const (
	PROTOCOL = "http"
)

type ClientFactory struct{}

func (f *ClientFactory) Create(address string) jsonwire.ClientInterface {
	return &Client{address: address}
}

func (c *Client) Address() string {
	return c.address
}

func (c *Client) Sessions() (*jsonwire.Sessions, error) {
	reqURL := url.URL{
		Scheme: PROTOCOL,
		Path:   "/wd/hub/sessions",
		Host:   c.Address(),
	}
	request, err := newRequest(http.MethodGet, reqURL.String(), "")
	if err != nil {
		return nil, fmt.Errorf("create json request, %v", err)
	}
	var sessions jsonwire.Sessions
	err = request.send(&sessions)
	if err != nil {
		return nil, fmt.Errorf("send json request, %v", err)
	}
	return &sessions, err
}

func (c *Client) Status() (*jsonwire.Message, error) {
	reqURL := url.URL{
		Scheme: PROTOCOL,
		Path:   "/wd/hub/status",
		Host:   c.Address(),
	}
	request, err := newRequest(http.MethodGet, reqURL.String(), "")
	if err != nil {
		return nil, fmt.Errorf("create json request, %v", err)
	}
	var message jsonwire.Message
	err = request.send(&message)
	if err != nil {
		return nil, fmt.Errorf("send json request, %v", err)
	}
	return &message, err
}

func (c *Client) CloseSession(sessionID string) (*jsonwire.Message, error) {
	reqURL := url.URL{
		Scheme: PROTOCOL,
		Path:   "/wd/hub/session/" + sessionID,
		Host:   c.Address(),
	}
	request, err := newRequest(http.MethodDelete, reqURL.String(), "")
	if err != nil {
		return nil, fmt.Errorf("create json request, %v", err)
	}
	var message jsonwire.Message
	err = request.send(&message)
	if err != nil {
		return nil, fmt.Errorf("send json request, %v", err)
	}
	return &message, nil
}

func newRequest(method, url string, requestBodyContent string) (*request, error) {
	b := strings.NewReader(requestBodyContent)
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, fmt.Errorf("create http request, %v", err)
	}
	if method == http.MethodPost {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-charset", "utf-8")
	return &request{httpRequest: req}, nil
}

// send as json.Unmarshal put result in variable pointed by outputStruct
func (req request) send(outputStruct interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second} //todo: move to config
	resp, err := client.Do(req.httpRequest)
	if err != nil {
		return fmt.Errorf("send http request, %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body, %v", err)
	}
	err = json.Unmarshal(body, outputStruct)
	if err != nil {
		return fmt.Errorf("unmarshal response error:[[%v]] message:[[%+v]]", err, outputStruct)
	}
	return nil
}
