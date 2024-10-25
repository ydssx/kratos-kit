package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type HttpClient interface {
	Get(url string, result interface{}) error
	Post(url string, payload interface{}, result interface{}) error
	Put(url string, payload interface{}, result interface{}) error
	Delete(url string, result interface{}) error
}

// Request represents a HTTP request
type Request struct {
	client      *resty.Client
	headers     map[string]string
	query       map[string]string
	formData    map[string]string
	timeout     time.Duration
	contentType string
	maxRetries  int
	retryDelay  time.Duration
}

// NewRequest creates a new Request instance with default settings
func NewRequest() *Request {
	return &Request{
		client:   resty.New().SetTimeout(30 * time.Second),
		headers:  make(map[string]string),
		query:    make(map[string]string),
		formData: make(map[string]string),
	}
}

func (r *Request) SetRetryOptions(maxRetries int, retryDelay time.Duration) *Request {
	r.maxRetries = maxRetries
	r.retryDelay = retryDelay
	r.client.SetRetryCount(r.maxRetries).SetRetryWaitTime(r.retryDelay)
	return r
}

// SetHeader sets the request header
func (r *Request) SetHeader(key, value string) *Request {
	r.headers[key] = value
	return r
}

// SetHeaders sets the request headers from the given map.
// The key-value pairs in the headers map will be set as the request headers.
func (r *Request) SetHeaders(headers map[string]string) *Request {
	for key, value := range headers {
		r.headers[key] = value
	}
	return r
}

// SetQuery sets the query parameters for the request
func (r *Request) SetQuery(key, value string) *Request {
	r.query[key] = value
	return r
}

// SetQuerys sets multiple query parameters from a map.
// The key-value pairs in the querys map will be set as the query parameters.
func (r *Request) SetQuerys(querys map[string]string) *Request {
	for k, v := range querys {
		r.SetQuery(k, v)
	}
	return r
}

// SetFormData sets the form data for the request
func (r *Request) SetFormData(key, value string) *Request {
	r.formData[key] = value
	return r
}

// SetTimeout sets the timeout duration for the request
func (r *Request) SetTimeout(timeout time.Duration) *Request {
	r.client.SetTimeout(timeout)
	return r
}

// WithContentType sets the content type of the request.
func (r *Request) WithContentType(contentType string) *Request {
	r.contentType = contentType
	return r
}

// Get 执行 GET 请求
func (r *Request) Get(url string, result interface{}) error {
	return r.doRequest("GET", url, nil, result)
}

// Post 执行 POST 请求
func (r *Request) Post(url string, payload interface{}, result interface{},) error {
	return r.doRequest("POST", url, payload, result)
}

// Put 执行 PUT 请求
func (r *Request) Put(url string, payload interface{}, result interface{}) error {
	return r.doRequest("PUT", url, payload, result)
}

// Delete sends a DELETE request to the specified URL with the given result interface{} and handles the response.
//
//   - url: The URL to send the DELETE request to.
//   - result: A pointer to the variable where the response will be stored.
//
// Returns:
//   - error: An error if there was a problem sending the request or handling the response.
func (r *Request) Delete(url string, result interface{}) error {
	return r.doRequest("DELETE", url, nil, result)
}

// doRequest 执行 HTTP 请求
func (r *Request) doRequest(method, url string, payload interface{}, result interface{}) error {
	req := r.client.R()
	r.addHeaders(req)
	r.addQueryParams(req)
	r.setContentType(req)

	if payload != nil {
		req.SetBody(payload)
	}

	resp, err := req.Execute(method, url)
	if err != nil {
		return err
	}

	return r.handleResponse(resp, result)
}

func (r *Request) addHeaders(req *resty.Request) {
	for key, value := range r.headers {
		req.SetHeader(key, value)
	}
}

func (r *Request) addQueryParams(req *resty.Request) {
	for key, value := range r.query {
		req.SetQueryParam(key, value)
	}
}

// setContentType 设置请求的 Content-Type 头
func (r *Request) setContentType(req *resty.Request) {
	if r.contentType != "" {
		req.SetHeader("Content-Type", r.contentType)
	}
}

// handleResponse 处理 HTTP 响应
func (r *Request) handleResponse(resp *resty.Response, result interface{}) error {
	if resp.StatusCode() >= http.StatusBadRequest {
		return fmt.Errorf("请求失败，状态码：%d，响应体：%s", resp.StatusCode(), resp.String())
	}

	if resp != nil && result != nil {
		return json.Unmarshal(resp.Body(), &result)
	}

	return nil
}
