package httpcheck

import (
	"encoding/json"
	"time"

	"github.com/premkit/healthcheck/log"
	"github.com/premkit/healthcheck/result"
)

// HTTPHealthcheck is a single node healthcheck on an endpoint which looks at the HTTP response code when querying.
type HTTPHealthcheck struct {
	Method                 string   `json:"method"`
	URI                    string   `json:"uri"`
	Options                []string `json:"options"`
	StatusCodesAvailable   []int    `json:"status_codes_available"`
	StatusCodesDegraded    []int    `json:"status_codes_unavailable"`
	StatusCodesUnavailable []int    `json:"status_codes_degraded"`
}

type HTTPHealthcheckResult struct {
	status      result.HealthcheckResponse
	TimeStarted time.Time
	Duration    time.Duration
	StatusCode  int
}

func (h *HTTPHealthcheck) Serialize() ([]byte, error) {
	b, err := json.Marshal(h)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return b, nil
}

func Deserialize(data []byte) (*HTTPHealthcheck, error) {
	httpCheck := HTTPHealthcheck{}
	if err := json.Unmarshal(data, &httpCheck); err != nil {
		log.Error(err)
		return nil, err
	}

	return &httpCheck, nil
}
