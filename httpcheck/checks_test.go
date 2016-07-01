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
