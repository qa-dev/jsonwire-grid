package selenium

import (
	"encoding/json"
	"errors"
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

type HttpMethod string

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
	reqUrl := url.URL{
		Scheme: PROTOCOL,
		Path:   "/wd/hub/sessions",
		Host:   c.Address(),
	}
	request, err := newRequest(http.MethodGet, reqUrl.String(), "")
	if err != nil {
		err = errors.New("Cant create request, " + err.Error())
		return nil, err
	}
	var sessions jsonwire.Sessions
	err = request.send(&sessions)
	if err != nil {
		err = errors.New("Cant read response, " + err.Error())
		return nil, err
	}
	return &sessions, err
}

func (c *Client) Status() (*jsonwire.Message, error) {
	reqUrl := url.URL{
		Scheme: PROTOCOL,
		Path:   "/wd/hub/status",
		Host:   c.Address(),
	}
	request, err := newRequest(http.MethodGet, reqUrl.String(), "")
	if err != nil {
		err = errors.New("Cant create request, " + err.Error())
		return nil, err
	}
	var message jsonwire.Message
	err = request.send(&message)
	if err != nil {
		err = errors.New("Cant read response, " + err.Error())
		return nil, err
	}
	return &message, err
}

func (c *Client) CloseSession(sessionId string) (*jsonwire.Message, error) {
	reqUrl := url.URL{
		Scheme: PROTOCOL,
		Path:   "/wd/hub/session/" + sessionId,
		Host:   c.Address(),
	}
	request, err := newRequest(http.MethodDelete, reqUrl.String(), "")
	if err != nil {
		err = errors.New("Cant create request, " + err.Error())
		return nil, err
	}
	var message jsonwire.Message
	err = request.send(&message)
	if err != nil {
		err = errors.New("Cant read response, " + err.Error())
		return nil, err
	}
	return &message, err
}

func newRequest(method, url string, requestBodyContent string) (*request, error) {
	b := strings.NewReader(requestBodyContent)
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, err
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
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req.httpRequest)
	if err != nil {
		err = errors.New("can't send request, " + err.Error())
		return err
	}
	defer resp.Body.Close()
	// todo: Получение респонза и разбор пока здесь.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("can't read response body, " + err.Error())
		return err
	}
	err = json.Unmarshal(body, outputStruct)
	if err != nil {
		err = errors.New("can't unmarshal response [" + string(body) + "], " + err.Error())
		return err
	}
	return nil
}
