package healthcheck

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Some mock types to make this testable.
type HealthcheckStepEmpty struct{}

func (h HealthcheckStepEmpty) RunSynchronously() HealthcheckStepResult {
	return HealthcheckStepResultEmpty{}
}

type HealthcheckStepResultEmpty struct{}

func (h HealthcheckStepResultEmpty) Status() HealthcheckResponse {
	return HealthcheckResponseAvailable
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

	o, err := CreateHealthcheck(&healthcheck)
	assert.Nil(t, o)
	assert.NotNil(t, err)

	multiError := err.(*multierror.Error)
	assert.Equal(t, 1, len(multiError.Errors))
	assert.Equal(t, errors.New("At least one step is required"), multiError.Errors[0])
}

func TestCreateHealthcheck(t *testing.T) {
	healthcheck := Healthcheck{
		Name:  "Empty",
		Steps: make([]HealthcheckStep, 0, 0),
	}
	healthcheck.Steps = append(healthcheck.Steps, HealthcheckStepEmpty{})

	o, err := CreateHealthcheck(&healthcheck)
	require.NoError(t, err)
	assert.Equal(t, 1, len(o.Steps), "there should be 1 step")
	assert.Equal(t, "Empty", o.Name, "name should be 'Empty'")
	assert.Equal(t, 64, len(o.ID), "id should be 64 chars long")
}
