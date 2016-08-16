package healthcheck

import (
	"errors"
	"testing"

	"github.com/premkit/healthcheck/result"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

// Some mock types to make this testable.
type HealthcheckStepEmpty struct{}

func (h HealthcheckStepEmpty) RunSynchronously() result.HealthcheckStepResult {
	return HealthcheckStepResultEmpty{}
}

type HealthcheckStepResultEmpty struct{}

func (h HealthcheckStepResultEmpty) Status() result.HealthcheckResponse {
	return result.HealthcheckResponseAvailable
}

// The tests
func TestValidateHealthcheckNoChecks(t *testing.T) {
	healthcheck := &Healthcheck{
		Name: "Name",
	}

	err := validateHealthcheck(healthcheck)
	assert.NotNil(t, err)

	multiError := err.(*multierror.Error)
	assert.Equal(t, 1, len(multiError.Errors))
	assert.Equal(t, errors.New("At least one step is required"), multiError.Errors[0])
}

func TestValidateHealthcheckNoNameNoChecks(t *testing.T) {
	healthcheck := &Healthcheck{}

	err := validateHealthcheck(healthcheck)
	assert.NotNil(t, err)

	multiError := err.(*multierror.Error)
	assert.Equal(t, 2, len(multiError.Errors))
	assert.Equal(t, errors.New("Name is required"), multiError.Errors[0])
	assert.Equal(t, errors.New("At least one step is required"), multiError.Errors[1])
}

func TestCreateHealthcheckEmpty(t *testing.T) {
	healthcheck := Healthcheck{
		Name: "Empty",
	}

	o, err := CreateOrUpdateHealthcheck(&healthcheck)
	assert.Nil(t, o)
	assert.NotNil(t, err)

	multiError := err.(*multierror.Error)
	assert.Equal(t, 1, len(multiError.Errors))
}
