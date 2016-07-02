package result

type HealthcheckResponse int

const (
	HealthcheckResponseUnknown HealthcheckResponse = iota
	HealthcheckResponseAvailable
	HealthcheckResponseDegraded
	HealthcheckResponseUnavailable
)

type HealthcheckStepResult interface {
	Status() HealthcheckResponse
}
