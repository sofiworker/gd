package ghttp

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"net"
	"net/url"
	"reflect"
	"time"
)

var (
	InvalidPathError    = fmt.Errorf("invalid path")
	BaseUrlEmptyError   = fmt.Errorf("baseurl is required when path is relative")
	BaseUrlFormatError  = fmt.Errorf("invalid baseurl")
	UrlNotAbsError      = fmt.Errorf("resulting url is not absolute")
	DataFormatError     = fmt.Errorf("data format error, only ptr data")
	NotFoundMethodError = fmt.Errorf("not found method")
)

type Client struct {
	baseUrl       string
	fastClient    *fasthttp.Client
	tlsConfig     *tls.Config
	defaultMethod string
}

type Request struct {
	fr             *fasthttp.Request
	client         *Client
	url            string
	requestBody    interface{}
	returnData     interface{}
	method         string
	streamBody     io.Reader
	streamBodySize int
}

type Response struct {
}

func NewClient() *Client {
	c := &fasthttp.Client{
		TLSConfig: &tls.Config{},
	}
	return &Client{fastClient: c}
}

func (c *Client) SetBaseUrl(baseUrl string) *Client {
	c.baseUrl = baseUrl
	return c
}

func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.fastClient.ReadTimeout = timeout
	c.fastClient.WriteTimeout = timeout
	return c
}

func (c *Client) SetDial(f func(addr string) (net.Conn, error)) *Client {
	c.fastClient.Dial = f
	return c
}

func (c *Client) SetTLSConfig(tlsConfig *tls.Config) *Client {
	c.tlsConfig = tlsConfig
	return c
}

func (c *Client) SetDefaultMethod(method string) *Client {
	c.defaultMethod = method
	return c
}

func (c *Client) SetBeforeRequestHook(...func(r *Request)) *Client {
	return c
}

func (c *Client) SeAfterRequestHook(...func(r *Request)) *Client {
	return c
}

func (c *Client) SetBeforeResponseHook() *Client {
	return c
}

func (c *Client) SetAfterResponseHook() *Client {
	return c
}

func (c *Client) SetCommonResponseBody() *Client {
	return c
}

func (c *Client) R() *Request {
	r := fasthttp.AcquireRequest()
	return &Request{fr: r, client: c}
}

func (r *Request) SetBearToken(token string) *Request {
	r.fr.Header.Add("Authorization", "Bearer "+token)
	return r
}

func (r *Request) SetJsonBody(data interface{}) *Request {
	r.requestBody = data
	return r
}

func (r *Request) SetUnmarshalData(data interface{}) *Request {
	r.returnData = data
	return r
}

func (r *Request) SetUrl(url string) *Request {
	r.url = url
	return r
}

func (r *Request) SetMethod(method string) *Request {
	r.method = method
	return r
}

func (r *Request) GetClient() *Client {
	return r.client
}

func (r *Request) GetFastHttpClient() *fasthttp.Client {
	return r.client.fastClient
}

func (r *Request) GetFastHttpRequest() *fasthttp.Request {
	return r.fr
}

func (r *Request) SetStreamBody(bodyStream io.Reader, bodySize int) *Request {
	r.streamBody = bodyStream
	r.streamBodySize = bodySize
	return r
}

func (r *Request) SetTracer() *Request {
	return r
}

func (r *Request) Done() (*Response, error) {
	if r.method == "" && r.client.defaultMethod == "" {
		return nil, NotFoundMethodError
	}

	if r.returnData != nil {
		if reflect.TypeOf(r.returnData).Kind() != reflect.Ptr {
			return nil, DataFormatError
		}
	}

	_, err := ConstructURL(r.client.baseUrl, r.url)
	if err != nil {
		return nil, err
	}

	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(r.fr)
	defer fasthttp.ReleaseResponse(response)

	bytes, err := json.Marshal(r.requestBody)
	if err != nil {
		return nil, err
	}
	r.fr.SetBody(bytes)

	err = r.client.fastClient.Do(r.fr, response)
	if err != nil {
		return nil, err
	}

	var resp *Response

	err = json.Unmarshal(response.Body(), r.returnData)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ConstructURL(baseurl, path string) (string, error) {
	pathURL, err := url.Parse(path)
	if err != nil {
		return "", InvalidPathError
	}

	if pathURL.IsAbs() {
		return pathURL.String(), nil
	}

	if baseurl == "" {
		return "", BaseUrlEmptyError
	}

	baseURL, err := url.Parse(baseurl)
	if err != nil {
		return "", BaseUrlFormatError
	}

	mergedURL := baseURL.ResolveReference(pathURL)

	if !mergedURL.IsAbs() {
		return "", UrlNotAbsError
	}

	return mergedURL.String(), nil
}
