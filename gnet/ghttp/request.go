package ghttp

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"reflect"
)

type Request struct {
	fr             *fasthttp.Request
	client         *Client
	url            string
	requestBody    interface{}
	returnData     interface{}
	method         string
	streamBody     io.Reader
	streamBodySize int
	enableDumpBody bool
	startRequest   int64
	endRequest     int64
	costRequest    int64
	startResponse  int64
	endResponse    int64
	costResponse   int64
	tracer         Tracer
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

func (r *Request) SetStreamBodyWithSize(bodyStream io.Reader, bodySize int) *Request {
	r.streamBody = bodyStream
	r.streamBodySize = bodySize
	return r
}

func (r *Request) SetStreamBody(bodyStream io.Reader) *Request {
	r.streamBody = bodyStream
	return r
}

func (r *Request) SetTracer() *Request {
	return r
}

func (r *Request) SetContentType(contentType string) *Request {
	r.fr.Header.Set("Content-Type", contentType)
	return r
}

func (r *Request) SetEnableDumpBody(enable bool) *Request {
	r.enableDumpBody = enable
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

	addr, err := ConstructURL(r.client.baseUrl, r.url)
	if err != nil {
		return nil, err
	}
	r.fr.SetRequestURI(addr)

	resp := new(Response)
	response := fasthttp.AcquireResponse()
	resp.fResp = response
	defer fasthttp.ReleaseRequest(r.fr)
	defer fasthttp.ReleaseResponse(resp.fResp)

	if r.requestBody != nil {
		bs, err := json.Marshal(r.requestBody)
		if err != nil {
			return nil, err
		}
		r.fr.SetBody(bs)
	}

	if r.tracer != nil {
		_, end := r.client.tracer.StartSpan()
		defer end()
	}
	err = r.client.fastClient.Do(r.fr, resp.fResp)
	if err != nil {
		return nil, err
	}
	resp.RemoteAddr = resp.fResp.RemoteAddr().String()
	resp.BodyRaw = resp.fResp.Body()

	if r.returnData != nil {
		err = json.Unmarshal(resp.fResp.Body(), r.returnData)
		if err != nil {
			return nil, err
		}
	}

	if r.enableDumpBody || r.client.enableDumpBody {
		fmt.Println(string(resp.fResp.Body()))
	}
	return resp, nil
}
