package httpcheck

import (
	"net/http"
	"testing"

	"github.com/premkit/healthcheck/healthcheck"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Run a http healthcheck test with no server available.
func TestHTTPHealthcheckRunNoServer(t *testing.T) {
	check := HTTPHealthcheck{
		Endpoint: "http://localhost/this-should-be-completely/invalid",
		Method:   "GET",
		ResponseAvailable: []int{
			http.StatusOK,
		},
	}

	resp, err := check.RunSynchronously()
	require.NoError(t, err)

	assert.Equal(t, healthcheck.HealthcheckResponseUnavailable, resp.Status())
}

// Run a http healthcheck test against google (passing)
func TestHTTPHealthcheckRunGetSuccess(t *testing.T) {
	check := HTTPHealthcheck{
		Endpoint: "http://google.com",
		Method:   "GET",
		ResponseAvailable: []int{
			http.StatusOK,
		},
	}

	resp, err := check.RunSynchronously()
	require.NoError(t, err)

	assert.Equal(t, healthcheck.HealthcheckResponseAvailable, resp.Status())
}

// Run a http healthcheck test that passes with a 404 (non-standard)
func TestHTTPHealthcheckRunGetSuccessOn404(t *testing.T) {
	check := HTTPHealthcheck{
		Endpoint: "https://github.com/premkit/invalid",
		Method:   "GET",
		ResponseAvailable: []int{
			http.StatusNotFound,
		},
	}

	resp, err := check.RunSynchronously()
	require.NoError(t, err)

	assert.Equal(t, healthcheck.HealthcheckResponseAvailable, resp.Status())
}

// Run a http healthcheck that returns a degraded response
func TestHTTPHealthcheckDegraded(t *testing.T) {
	check := HTTPHealthcheck{
		Endpoint: "https://github.com/premkit/invalid",
		Method:   "GET",
		ResponseAvailable: []int{
			http.StatusOK,
		},
		ResponseDegraded: []int{
			http.StatusNotFound,
		},
		ResponseUnavailable: []int{
			http.StatusCreated,
		},
	}

	resp, err := check.RunSynchronously()
	require.NoError(t, err)

	assert.Equal(t, healthcheck.HealthcheckResponseDegraded, resp.Status())
}

// Run an http healthcheck that returns an unavailable response
func TestHTTPHealthcheckUnavailable(t *testing.T) {
	check := HTTPHealthcheck{
		Endpoint: "https://github.com/premkit/invalid",
		Method:   "GET",
		ResponseAvailable: []int{
			http.StatusOK,
		},
		ResponseDegraded: []int{
			http.StatusCreated,
		},
		ResponseUnavailable: []int{
			http.StatusNotFound,
			http.StatusForbidden,
		},
	}

	resp, err := check.RunSynchronously()
	require.NoError(t, err)

	assert.Equal(t, healthcheck.HealthcheckResponseUnavailable, resp.Status())
}

// Run an http healthcheck with an unknown method
func TestHTTPHealthcheckUnsupported(t *testing.T) {
	check := HTTPHealthcheck{
		Endpoint: "https://github.com/premkit/invalid",
		Method:   "OPTIONS",
		ResponseAvailable: []int{
			http.StatusOK,
		},
	}

	resp, err := check.RunSynchronously()
	assert.Equal(t, "Unsupported method \"OPTIONS\"", err.Error())
	assert.Equal(t, healthcheck.HealthcheckResponseUnknown, resp.Status())
}

// Run an http healthcheck with a POST method
func TestHTTPHealthcheckPost(t *testing.T) {
	check := HTTPHealthcheck{
		Endpoint: "https://github.com/premkit/invalid",
		Method:   "POST",
		ResponseUnavailable: []int{
			http.StatusNotFound,
		},
	}

	_, err := check.RunSynchronously()
	assert.Equal(t, "Not implemented", err.Error())
}

// Run an http healthcheck with a PUT method
func TestHTTPHealthcheckPut(t *testing.T) {
	check := HTTPHealthcheck{
		Endpoint: "https://github.com/premkit/invalid",
		Method:   "PUT",
		ResponseUnavailable: []int{
			http.StatusNotFound,
		},
	}

	_, err := check.RunSynchronously()
	assert.Equal(t, "Not implemented", err.Error())
}
