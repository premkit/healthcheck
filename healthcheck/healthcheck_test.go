package healthcheck

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

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

func TestCreateHealthcheck_Empty(t *testing.T) {
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
