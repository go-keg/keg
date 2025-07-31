package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// setup starts a test server and returns its URL and a teardown function.
func setup(t *testing.T, handler http.HandlerFunc) (string, func()) {
	server := httptest.NewServer(handler)
	return server.URL, func() {
		server.Close()
	}
}

func TestNewClient(t *testing.T) {
	t.Run("WithTimeout", func(t *testing.T) {
		timeout := 10 * time.Second
		client := NewClient(WithTimeout(timeout))
		assert.Equal(t, timeout, client.client.Timeout)
	})

	t.Run("WithBaseURL", func(t *testing.T) {
		baseURL := "http://example.com/api"
		client := NewClient(WithBaseURL(baseURL))
		assert.Equal(t, baseURL, client.baseURL)
	})

	t.Run("WithProxy", func(t *testing.T) {
		proxyURL := "http://proxy.example.com:8080"
		client := NewClient(WithProxy(proxyURL))
		transport, ok := client.client.Transport.(*http.Transport)
		assert.True(t, ok)
		proxy, err := transport.Proxy(&http.Request{})
		assert.NoError(t, err)
		assert.Equal(t, proxyURL, proxy.String())
	})
}

func TestClient_Request(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			assert.Equal(t, "bar", r.URL.Query().Get("foo"))
			assert.Equal(t, "my-header-value", r.Header.Get("X-My-Header"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"message": "get successful"}`)
		case http.MethodPost:
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			body, _ := io.ReadAll(r.Body)
			var data map[string]string
			_ = json.Unmarshal(body, &data)
			assert.Equal(t, "world", data["hello"])
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, `{"message": "post successful"}`)
		case http.MethodPut:
			assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
			r.ParseForm()
			assert.Equal(t, "bar", r.FormValue("foo"))
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"message": "put successful"}`)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

	serverURL, teardown := setup(t, handler)
	defer teardown()

	client := NewClient()

	t.Run("GET", func(t *testing.T) {
		resp, err := client.Get(context.Background(), serverURL,
			SetQueryParams(map[string]string{"foo": "bar"}),
			SetHeader("X-My-Header", "my-header-value"),
		)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		assert.Equal(t, "200 OK", resp.Status())
		assert.Equal(t, "application/json", resp.GetHeader("Content-Type"))
		var result map[string]string
		err = resp.Unmarshal(&result)
		assert.NoError(t, err)
		assert.Equal(t, "get successful", result["message"])
	})

	t.Run("POST", func(t *testing.T) {
		resp, err := client.Post(context.Background(), serverURL, SetBody(map[string]string{"hello": "world"}))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode())
		var result map[string]string
		err = resp.Unmarshal(&result)
		assert.NoError(t, err)
		assert.Equal(t, "post successful", result["message"])
	})

	t.Run("PUT with FormData", func(t *testing.T) {
		resp, err := client.Put(context.Background(), serverURL, SetFormData(map[string]string{"foo": "bar"}))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		var result map[string]string
		err = resp.Unmarshal(&result)
		assert.NoError(t, err)
		assert.Equal(t, "put successful", result["message"])
	})

	t.Run("DELETE", func(t *testing.T) {
		resp, err := client.Delete(context.Background(), serverURL)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode())
		assert.Empty(t, resp.Content())
	})
}

func TestClient_BaseURL(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/users", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}

	serverURL, teardown := setup(t, handler)
	defer teardown()

	client := NewClient(WithBaseURL(serverURL + "/api/v1"))
	resp, err := client.Get(context.Background(), "/users")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestClient_Graphql(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		var body map[string]any
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, "query { users { id name } }", body["query"])
		vars, ok := body["variables"].(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "test", vars["role"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"data":{"users":[{"id":"1","name":"test"}]}}`)
	}

	serverURL, teardown := setup(t, handler)
	defer teardown()

	client := NewClient()
	query := "query { users { id name } }"
	variables := map[string]any{"role": "test"}

	resp, err := client.Graphql(context.Background(), serverURL, query, variables)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	var result map[string]any
	err = resp.Unmarshal(&result)
	assert.NoError(t, err)
	assert.NotNil(t, result["data"])
}

func TestSetPathParams(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users/123/posts/abc", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}

	serverURL, teardown := setup(t, handler)
	defer teardown()

	client := NewClient()
	resp, err := client.Get(context.Background(), serverURL+"/users/{userID}/posts/{postID}", SetPathParams(map[string]string{
		"userID": "123",
		"postID": "abc",
	}))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestResponse(t *testing.T) {
	header := http.Header{}
	header.Set("X-Test", "value")
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	resp := Response{
		status:     "200 OK",
		statusCode: 200,
		content:    []byte(`{"key":"val"}`),
		headers:    header,
		request:    req,
	}

	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Equal(t, `{"key":"val"}`, resp.Content())
	assert.Equal(t, "value", resp.GetHeader("X-Test"))
	assert.Equal(t, header, resp.Header())
	assert.Equal(t, req, resp.Request())

	var data map[string]string
	err := resp.Unmarshal(&data)
	assert.NoError(t, err)
	assert.Equal(t, "val", data["key"])
}