package client

import (
	"io"
	"net/http"
)

type Clienter interface {
	Request(t, uri string, headers, params map[string]string, body io.Reader) (*http.Response, error)
}

type Client struct {
	client http.Client
}

func NewClient(client http.Client) *Client {
	return &Client{client: client}
}

func (c *Client) Request(t, uri string, headers, params map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(t, uri, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	return c.client.Do(req)
}
