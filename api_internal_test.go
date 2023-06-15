package portfolio

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func Test_portfolioTreeURL_default(t *testing.T) {
	t.Setenv(ServerURLEnvironmentVariableName, "")
	assert.Equal(t, DefaultURL, portfolioTreeURL())
}

func Test_portfolioTreeURL_other(t *testing.T) {
	t.Setenv(ServerURLEnvironmentVariableName, "other")
	assert.Equal(t, "other", portfolioTreeURL())
}

func Test_doJSONRequest_do_fails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := doJSONRequest[struct{}](func(r *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("lemon")
	}, req)
	assert.Error(t, err)
}

func Test_doJSONRequest_unexpected_status(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := doJSONRequest[struct{}](func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusTeapot,
			Body:       io.NopCloser(&bytes.Reader{}),
		}, nil
	}, req)
	assert.Error(t, err)
}

func Test_doJSONRequest_read_fails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := doJSONRequest[struct{}](func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(iotest.ErrReader(fmt.Errorf("lemon"))),
		}, nil
	}, req)
	assert.Error(t, err)
}

func Test_doJSONRequest_invalid_json(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := doJSONRequest[struct{}](func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("[]")),
		}, nil
	}, req)
	assert.Error(t, err)
}
