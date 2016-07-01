package schema

import (
	"github.com/premkit/healthcheck/healthcheck"
)

type CreateHealthcheckRequest struct {
	Healthcheck *healthcheck.Healthcheck `json:"healthcheck"`
}

type CreateHealthcheckResponse struct {
	Healthcheck *healthcheck.Healthcheck `json:"healthcheck"`
}

type ListHealthchecksResponse struct {
	Healthchecks []*healthcheck.Healthcheck `json:"healthchecks"`
}
