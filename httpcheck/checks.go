package httpcheck

import (
	"net/http"
	"time"

	"github.com/premkit/healthcheck/result"

	"github.com/parnurzeal/gorequest"
)

func (r HTTPHealthcheckResult) Status() result.HealthcheckResponse {
	return r.status
}

// RunSynchronously will run the test and return the results immediately.
func (h HTTPHealthcheck) RunSynchronously() result.HealthcheckStepResult {
	if h.Method == "GET" {
		return h.getSynchronously()
	} else if h.Method == "POST" {
		return h.postSynchronously()
	} else if h.Method == "PUT" {
		return h.putSynchronously()
	}

	// Use an unknown response here.
	response := HTTPHealthcheckResult{
		status:      result.HealthcheckResponseUnknown,
		TimeStarted: time.Now(),
	}

	return result.HealthcheckStepResult(response)
}

func (h HTTPHealthcheck) getSynchronously() HTTPHealthcheckResult {
	response := HTTPHealthcheckResult{
		status:      result.HealthcheckResponseUnknown,
		TimeStarted: time.Now(),
	}

	// Make the request
	resp, _, _ := gorequest.New().Get(h.URI).End()

	response.Duration = time.Now().Sub(response.TimeStarted)
	if resp == nil {
		response.status = result.HealthcheckResponseUnavailable
	} else {
		response.status = result.HealthcheckResponseUnavailable
		for _, r := range h.StatusCodesAvailable {
			if resp.StatusCode == r {
				response.status = result.HealthcheckResponseAvailable
			}
		}
		for _, r := range h.StatusCodesDegraded {
			if resp.StatusCode == r {
				response.status = result.HealthcheckResponseDegraded
			}
		}
		for _, r := range h.StatusCodesUnavailable {
			if resp.StatusCode == r {
				response.status = result.HealthcheckResponseUnavailable
			}
		}
		response.StatusCode = http.StatusInternalServerError
	}

	return response
}

func (h HTTPHealthcheck) postSynchronously() HTTPHealthcheckResult {
	response := HTTPHealthcheckResult{
		status:      result.HealthcheckResponseUnknown,
		TimeStarted: time.Now(),
	}

	return response
}

func (h HTTPHealthcheck) putSynchronously() HTTPHealthcheckResult {
	response := HTTPHealthcheckResult{
		status:      result.HealthcheckResponseUnknown,
		TimeStarted: time.Now(),
	}

	return response
}
