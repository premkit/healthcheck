package httpcheck

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/premkit/healthcheck/healthcheck"

	"github.com/parnurzeal/gorequest"
)

func (r *HTTPHealthcheckResult) Status() healthcheck.HealthcheckResponse {
	return r.status
}

// RunSynchronously will run the test and return the results immediately.
func (h *HTTPHealthcheck) RunSynchronously() (*HTTPHealthcheckResult, error) {
	if h.Method == "GET" {
		return h.getSynchronously()
	} else if h.Method == "POST" {
		return h.postSynchronously()
	} else if h.Method == "PUT" {
		return h.putSynchronously()
	}

	return nil, fmt.Errorf("Unsupported method %q", h.Method)
}

func (h *HTTPHealthcheck) getSynchronously() (*HTTPHealthcheckResult, error) {
	response := &HTTPHealthcheckResult{
		status:      healthcheck.HealthcheckResponseUnknown,
		TimeStarted: time.Now(),
	}

	// Make the request
	resp, _, _ := gorequest.New().Get(h.Endpoint).End()

	response.Duration = time.Now().Sub(response.TimeStarted)
	if resp == nil {
		response.status = healthcheck.HealthcheckResponseUnavailable
	} else {
		response.status = healthcheck.HealthcheckResponseUnavailable
		for _, r := range h.ResponseAvailable {
			if resp.StatusCode == r {
				response.status = healthcheck.HealthcheckResponseAvailable
			}
		}
		for _, r := range h.ResponseDegraded {
			if resp.StatusCode == r {
				response.status = healthcheck.HealthcheckResponseDegraded
			}
		}
		for _, r := range h.ResponseUnavailable {
			if resp.StatusCode == r {
				response.status = healthcheck.HealthcheckResponseUnavailable
			}
		}
		response.ResponseCode = http.StatusInternalServerError
	}

	return response, nil
}

func (h *HTTPHealthcheck) postSynchronously() (*HTTPHealthcheckResult, error) {
	return nil, errors.New("Not implemented")
}

func (h *HTTPHealthcheck) putSynchronously() (*HTTPHealthcheckResult, error) {
	return nil, errors.New("Not implemented")
}
