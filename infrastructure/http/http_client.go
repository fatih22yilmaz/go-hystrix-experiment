package http

import (
	"github.com/valyala/fasthttp"
	"time"
)

const AgentName string = "hystrix-experiment-client"

type httpClient struct {
	client *fasthttp.Client
}

type Client interface {
	Get(host string, queryString string) ([]byte, error)
}

func NewHTTPClient() Client {
	return &httpClient{
		&fasthttp.Client{
			Name:        AgentName,
			ReadTimeout: time.Duration(750) * time.Millisecond,
		},
	}
}

func (httpClient *httpClient) Get(host string, queryString string) ([]byte, error) {
	res := fasthttp.AcquireResponse()
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseResponse(res)
	defer fasthttp.ReleaseRequest(req)

	url := host + queryString

	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("x-agentName", AgentName)

	err := httpClient.client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return res.Body(), nil
}
