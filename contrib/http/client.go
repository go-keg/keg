package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type Client struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

type ClientOptionFunc func(c *Client)

func WithProxy(u string) ClientOptionFunc {
	urlParse, err := neturl.Parse(u)
	if err != nil {
		return func(c *Client) {}
	}
	return func(c *Client) {
		c.client.Transport = &http.Transport{
			Proxy: http.ProxyURL(urlParse),
		}
	}
}

func WithTimeout(d time.Duration) ClientOptionFunc {
	return func(c *Client) {
		c.client.Timeout = d
	}
}

func WithBaseURL(baseURL string) ClientOptionFunc {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithHeaders(headers map[string]string) ClientOptionFunc {
	return func(c *Client) {
		c.headers = headers
	}
}

func NewClient(opts ...ClientOptionFunc) *Client {
	c := Client{client: &http.Client{}}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

type Response struct {
	status     string
	statusCode int
	content    []byte
	headers    http.Header
	request    *http.Request
}

type OptionFunc func(r *http.Request)

func (c Client) Client() *http.Client {
	return c.client
}

func (c Client) Get(ctx context.Context, url string, opts ...OptionFunc) (*Response, error) {
	return c.Request(ctx, http.MethodGet, url, opts...)
}

func (c Client) Post(ctx context.Context, url string, opts ...OptionFunc) (*Response, error) {
	return c.Request(ctx, http.MethodPost, url, opts...)
}

func (c Client) Put(ctx context.Context, url string, opts ...OptionFunc) (*Response, error) {
	return c.Request(ctx, http.MethodPut, url, opts...)
}

func (c Client) Delete(ctx context.Context, url string, opts ...OptionFunc) (*Response, error) {
	return c.Request(ctx, http.MethodDelete, url, opts...)
}

func (c Client) Graphql(ctx context.Context, endpoint, query string, variables map[string]any, opts ...OptionFunc) (*Response, error) {
	return c.Post(ctx, endpoint, append(opts, SetBody(map[string]any{
		"query":     query,
		"variables": variables,
	}))...)
}

func (c Client) Request(ctx context.Context, method, url string, opts ...OptionFunc) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(req)
	}
	if c.baseURL != "" {
		parse, err := neturl.Parse(c.baseURL)
		if err != nil {
			return nil, err
		}
		if req.URL.Host == "" {
			req.URL.Host = parse.Host
		}
		if req.URL.Scheme == "" {
			req.URL.Scheme = parse.Scheme
		}
		req.URL.Path = parse.JoinPath(req.URL.Path).Path
	}
	for k, v := range c.headers {
		if req.Header.Get(k) == "" {
			req.Header.Set(k, v)
		}
	}
	response, err := c.client.Do(req)
	if err != nil {
		return &Response{
			headers: req.Header,
			request: req,
		}, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	return c.response(response)
}

func SetHeaders(headers map[string]string) OptionFunc {
	return func(r *http.Request) {
		for key, value := range headers {
			r.Header.Add(key, value)
		}
	}
}

func SetHeader(key, value string) OptionFunc {
	return func(r *http.Request) {
		r.Header.Add(key, value)
	}
}

func SetPathParams(params map[string]string) OptionFunc {
	return func(r *http.Request) {
		for k, v := range params {
			r.URL.Path = strings.ReplaceAll(r.URL.Path, "{"+k+"}", neturl.PathEscape(v))
		}
	}
}

func SetFormData(data map[string]any) OptionFunc {
	return func(r *http.Request) {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Body = io.NopCloser(strings.NewReader(map2Values(data).Encode()))
	}
}

func SetBody(data any) OptionFunc {
	return func(r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		marshal, _ := json.Marshal(data)
		r.Body = io.NopCloser(bytes.NewBuffer(marshal))
	}
}

func map2Values(params map[string]any) neturl.Values {
	query := neturl.Values{}
	for key, val := range params {
		switch v := val.(type) {
		case int, int64, float64, bool:
			query.Set(key, cast.ToString(v))
		case string:
			query.Set(key, v)
		case []string:
			query.Del(key)
			for _, item := range v {
				query.Add(key, item)
			}
		case []int:
			query.Del(key)
			for _, item := range v {
				query.Add(key, cast.ToString(item))
			}
		case []int64:
			query.Del(key)
			for _, item := range v {
				query.Add(key, cast.ToString(item))
			}
		case []float64:
			query.Del(key)
			for _, item := range v {
				query.Add(key, cast.ToString(item))
			}
		default:
			log.Warn(fmt.Sprintf("unsupported type: %T, key: %s", v, key))
		}
	}
	return query
}

func SetQueryParams(params map[string]any) OptionFunc {
	return func(r *http.Request) {
		r.URL.RawQuery = map2Values(params).Encode()
	}
}

func (c Client) response(resp *http.Response) (*Response, error) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{
		status:     resp.Status,
		statusCode: resp.StatusCode,
		content:    body,
		headers:    resp.Header,
		request:    resp.Request,
	}, err
}

func (r Response) StatusCode() int {
	return r.statusCode
}

func (r Response) Status() string {
	return r.status
}

func (r Response) Content() string {
	return string(r.content)
}

func (r Response) Unmarshal(v any) error {
	return json.Unmarshal(r.content, v)
}

func (r Response) GetHeader(key string) string {
	return r.headers.Get(key)
}

func (r Response) Header() http.Header {
	return r.headers
}

func (r Response) Request() *http.Request {
	return r.request
}
