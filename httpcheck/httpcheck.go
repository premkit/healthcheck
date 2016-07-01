package httpcheck

import (
	"time"

	"github.com/premkit/healthcheck/healthcheck"
)

// HTTPHealthcheck is a single node healthcheck on an endpoint which looks at the HTTP response code when querying.
type HTTPHealthcheck struct {
	Endpoint            string                 `json:"endpoint"`
	Method              string                 `json:"method"`
	Body                string                 `json:"body"`
	Headers             map[string]interface{} `json:"headers"`
	FollowRedirects     bool                   `json:"follow_redirects"`
	AllowInsecure       bool                   `json:"allow_insecure"`
	ResponseAvailable   []int                  `json:"response_available"`
	ResponseDegraded    []int                  `json:"response_degraded"`
	ResponseUnavailable []int                  `json:"response_unavailable"`
}

type HTTPHealthcheckResult struct {
	status       healthcheck.HealthcheckResponse
	TimeStarted  time.Time
	Duration     time.Duration
	ResponseCode int
}
