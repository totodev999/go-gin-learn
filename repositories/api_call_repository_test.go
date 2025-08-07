package repositories

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

// モックRoundTripper
type mockRoundTripper struct {
	fn func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.fn(req)
}

// ポインタユーティリティ
func strPtr(s string) *string { return &s }
func uintPtr(i uint) *uint    { return &i }
func intPtr(i int) *int       { return &i }

func TestGetAllPosts_Success(t *testing.T) {
	t.Setenv("BASE_URL", "http://dummy")

	body := `[{"id":1,"title":"hello","body":"body","userId":2,"dummy":"dummy"}]`
	mockRT := &mockRoundTripper{
		fn: func(req *http.Request) (*http.Response, error) {
			t.Logf("\nMock RoundTrip called. URL.Path=%s\n", req.URL)
			if req.URL.Path == "/posts" {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Request:    req,
				}, nil
			}
			return nil, fmt.Errorf("unexpected path %s", req.URL.Path)
		},
	}

	client := resty.New()
	client.SetTransport(mockRT)
	repo := NewAPICallRepository(client)

	got, err := repo.GetAllPosts(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, uint(1), (*(*got)[0].Id))
	assert.Equal(t, "hello", (*(*got)[0].Title))
}

func TestGetAllPosts_ErrorStatus(t *testing.T) {
	t.Setenv("BASE_URL", "http://dummy")
	mockRT := &mockRoundTripper{
		fn: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader("error")),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Request:    req,
			}, nil
		},
	}
	client := resty.New()
	client.SetTransport(mockRT)
	repo := NewAPICallRepository(client)

	got, err := repo.GetAllPosts(context.Background())
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestGetUserAndPosts_Success(t *testing.T) {
	t.Setenv("BASE_URL", "http://dummy")

	userBody := `{"id":1,"name":"foo","username":"bar","email":"baz","address":{},"phone":"","website":"","company":{}}`
	postsBody := `[{"id":1,"title":"test title","body":"test body","userId":1,"dummy":null}]`

	mockRT := &mockRoundTripper{
		fn: func(req *http.Request) (*http.Response, error) {
			switch {
			case req.URL.Path == "/users/1":
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(userBody)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Request:    req,
				}, nil
			case req.URL.Path == "/posts" && req.URL.Query().Get("userId") == "1":
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(postsBody)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Request:    req,
				}, nil
			default:
				return &http.Response{
					StatusCode: 404,
					Body:       io.NopCloser(strings.NewReader("not found")),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Request:    req,
				}, nil
			}
		},
	}
	client := resty.New()
	client.SetTransport(mockRT)
	repo := NewAPICallRepository(client)

	got, err := repo.GetUserAndPosts(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, 1, got.User.ID)
	assert.Equal(t, "foo", got.User.Name)
	assert.Equal(t, "test title", *got.Posts[0].Title)
}

func TestGetUserAndPosts_UserApiError(t *testing.T) {
	t.Setenv("BASE_URL", "http://dummy")
	mockRT := &mockRoundTripper{
		fn: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path == "/users/1" {
				return &http.Response{
					StatusCode: 500,
					Body:       io.NopCloser(strings.NewReader("error")),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Request:    req,
				}, nil
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("[]")),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Request:    req,
			}, nil
		},
	}
	client := resty.New()
	client.SetTransport(mockRT)
	repo := NewAPICallRepository(client)

	got, err := repo.GetUserAndPosts(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, got)
}
