package schema

import (
	"github.com/premkit/healthcheck/healthcheck"
)

type ListHealthchecksResponse struct {
	Healthchecks []*healthcheck.Healthcheck `json:"healthchecks"`
}

type RunHealthchecksResponse struct {
}
