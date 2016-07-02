package httpcheck

import (
	"net/http"
	"testing"

	"github.com/premkit/healthcheck/result"

	"github.com/stretchr/testify/assert"
)

// Run a http healthcheck test with no server available.
func TestHTTPHealthcheckRunNoServer(t *testing.T) {
	check := HTTPHealthcheck{
		URI:    "http://localhost/this-should-be-completely/invalid",
		Method: "GET",
		StatusCodesAvailable: []int{
			http.StatusOK,
		},
	}

	resp := check.RunSynchronously()
	assert.Equal(t, result.HealthcheckResponseUnavailable, resp.Status())
}

// Run a http healthcheck test against google (passing)
func TestHTTPHealthcheckRunGetSuccess(t *testing.T) {
	check := HTTPHealthcheck{
		URI:    "http://google.com",
		Method: "GET",
		StatusCodesAvailable: []int{
			http.StatusOK,
		},
	}

	resp := check.RunSynchronously()
	assert.Equal(t, result.HealthcheckResponseAvailable, resp.Status())
}

// Run a http healthcheck test that passes with a 404 (non-standard)
func TestHTTPHealthcheckRunGetSuccessOn404(t *testing.T) {
	check := HTTPHealthcheck{
		URI:    "https://github.com/premkit/invalid",
		Method: "GET",
		StatusCodesAvailable: []int{
			http.StatusNotFound,
		},
	}

	resp := check.RunSynchronously()
	assert.Equal(t, result.HealthcheckResponseAvailable, resp.Status())
}

// Run a http healthcheck that returns a degraded response
func TestHTTPHealthcheckDegraded(t *testing.T) {
	check := HTTPHealthcheck{
		URI:    "https://github.com/premkit/invalid",
		Method: "GET",
		StatusCodesAvailable: []int{
			http.StatusOK,
		},
		StatusCodesDegraded: []int{
			http.StatusNotFound,
		},
		StatusCodesUnavailable: []int{
			http.StatusCreated,
		},
	}

	resp := check.RunSynchronously()
	assert.Equal(t, result.HealthcheckResponseDegraded, resp.Status())
}

// Run an http healthcheck that returns an unavailable response
func TestHTTPHealthcheckUnavailable(t *testing.T) {
	check := HTTPHealthcheck{
		URI:    "https://github.com/premkit/invalid",
		Method: "GET",
		StatusCodesAvailable: []int{
			http.StatusOK,
		},
		StatusCodesDegraded: []int{
			http.StatusCreated,
		},
		StatusCodesUnavailable: []int{
			http.StatusNotFound,
			http.StatusForbidden,
		},
	}

	resp := check.RunSynchronously()
	assert.Equal(t, result.HealthcheckResponseUnavailable, resp.Status())
}

// Run an http healthcheck with an unknown method
func TestHTTPHealthcheckUnsupported(t *testing.T) {
	check := HTTPHealthcheck{
		URI:    "https://github.com/premkit/invalid",
		Method: "OPTIONS",
		StatusCodesAvailable: []int{
			http.StatusOK,
		},
	}

	resp := check.RunSynchronously()
	assert.Equal(t, result.HealthcheckResponseUnknown, resp.Status())
}

// Run an http healthcheck with a POST method
func TestHTTPHealthcheckPost(t *testing.T) {
	/*	check := HTTPHealthcheck{
			Endpoint: "https://github.com/premkit/invalid",
			Method:   "POST",
			ResponseUnavailable: []int{
				http.StatusNotFound,
			},
		}

		_, _ = check.RunSynchronously()*/
}

// Run an http healthcheck with a PUT method
func TestHTTPHealthcheckPut(t *testing.T) {
	/*	check := HTTPHealthcheck{
			Endpoint: "https://github.com/premkit/invalid",
			Method:   "PUT",
			ResponseUnavailable: []int{
				http.StatusNotFound,
			},
		}

		_, err := check.RunSynchronously()
		assert.Equal(t, "Not implemented", err.Error())*/
}
