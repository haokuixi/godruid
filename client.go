package godruid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	DefaultEndPoint = "/druid/v2"
)

type Client struct {
	Url      string
	EndPoint string
	Username string
	Password string

	Debug        bool
	LastRequest  string
	LastResponse string
	HttpClient   *http.Client
}

func (c *Client) Query(query Query) (err error) {
	query.setup()
	var reqJson []byte
	if c.Debug {
		reqJson, err = json.MarshalIndent(query, "", "  ")
	} else {
		reqJson, err = json.Marshal(query)
	}
	if err != nil {
		return
	}

	result, err := c.QueryRaw(reqJson)
	if err != nil {
		return
	}

	return query.onResponse(result)
}

func (c *Client) QueryRaw(req []byte) (result []byte, err error) {
	if c.EndPoint == "" {
		c.EndPoint = DefaultEndPoint
	}
	endPoint := c.EndPoint
	if c.Debug {
		endPoint += "?pretty"
		c.LastRequest = string(req)
	}
	if err != nil {
		return
	}

	request, err := http.NewRequest("POST", c.Url+endPoint, bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	if c.Username != "" {
		request.SetBasicAuth(c.Username, c.Password)
	}
	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if c.Debug {
		c.LastResponse = string(result)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(result))
	}

	return
}
