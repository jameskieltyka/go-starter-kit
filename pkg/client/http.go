package client

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"go.uber.org/zap"
)

type HTTPClient struct {
	*http.Client
	middleware []HTTPClientMiddleware
	method     string
	url        string
	body       io.Reader
}

type HTTPClientMiddleware interface {
	Run(ctx context.Context, f ExecutableOperation) (*http.Response, error)
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		middleware: make([]HTTPClientMiddleware, 0),
	}
}

func (h *HTTPClient) WithMiddleware(middleware ...HTTPClientMiddleware) *HTTPClient {
	h.middleware = append(h.middleware, middleware...)
	return h
}

func (h *HTTPClient) WithRequest(method string, url string, body io.Reader) *HTTPClient {
	h.method = method
	h.url = url
	h.body = body
	return h
}

func (h *HTTPClient) Run(ctx context.Context) (*http.Response, error) {
	operation := func(ctx context.Context) (*http.Response, error) {
		req, err := http.NewRequest(h.method, h.url, h.body)
		if err != nil {
			return nil, err
		}
		client := http.DefaultClient
		return client.Do(req)
	}

	for i := range h.middleware {
		operation = WrapMiddleware(h.middleware[i], operation)
	}

	return operation(ctx)
}

func WrapMiddleware(f HTTPClientMiddleware, p ExecutableOperation) ExecutableOperation {
	return func(ctx context.Context) (*http.Response, error) {
		return f.Run(ctx, p)
	}
}

type ExecutableOperation func(ctx context.Context) (*http.Response, error)

type CircuitBreaker struct {
	Name string
}

func (c CircuitBreaker) Configure() HTTPClientMiddleware {
	hystrix.ConfigureCommand(c.Name, hystrix.CommandConfig{
		Timeout:               1000,
		MaxConcurrentRequests: 1,
		ErrorPercentThreshold: 25,
	})
	return c
}

func (c CircuitBreaker) Run(ctx context.Context, f ExecutableOperation) (*http.Response, error) {
	fmt.Println("running circuit breaker")
	var result *http.Response
	var err error

	err = hystrix.Do(c.Name, func() error {
		result, err = f(ctx)
		return err
	}, nil)
	return result, err
}

type Logger struct {
	Url    string
	Method string
}

func (l Logger) Run(ctx context.Context, f ExecutableOperation) (*http.Response, error) {
	fmt.Println("running logger")
	zap.S().Infof("request: method %s, uri %s", l.Method, l.Url)
	result, err := f(ctx)
	zap.S().Infof("response: result %v, err %v", result, err)
	return result, err
}
