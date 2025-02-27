package ghttp

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net"
	"time"
)

type Client struct {
	baseUrl string
	client  *fasthttp.Client
}

type Request struct {
	fr *fasthttp.Request
	c  *Client
}

func NewClient() *Client {
	c := &fasthttp.Client{}
	return &Client{client: c}
}

func (c *Client) SetBaseUrl(baseUrl string) {
	c.baseUrl = baseUrl
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.ReadTimeout = timeout
	c.client.WriteTimeout = timeout
}

func (c *Client) SetDial(f func(addr string) (net.Conn, error)) {
	c.client.Dial = f
}

func (c *Client) SetBear() {
	//c.client.
}

func (c *Client) R() *Request {
	r := fasthttp.AcquireRequest()
	return &Request{fr: r, c: c}
}

func (r *Request) SetBearToken() {
}

func (r *Request) SetJsonBody(data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.fr.SetBody(bytes)
	return nil
}
