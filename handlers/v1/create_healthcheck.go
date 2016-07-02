package v1

import (
	"net/http"

	"github.com/premkit/healthcheck/healthcheck"
	"github.com/premkit/healthcheck/log"

	"github.com/gin-gonic/gin"
)

// CreateHealthcheckParams contains parameters to the create healthcheck route.
// swagger:parameters createHealthcheck
type CreateHealthcheckParams struct {
	// Healthcheck registration parameters.
	// In: body
	Healthcheck *healthcheck.Healthcheck `json:"healthcheck"`
}

// CreateHealthcheckResponse represents the response to a createHealthcheck call. This response
// includes a pointer to the created healthcheck
// swagger:response createHealthcheck
type CreateHealthcheckResponse struct {
	// Healthcheck
	// In: body
	Healthcheck *healthcheck.Healthcheck `json:"healthcheck"`
}

// CreateHealthcheck is called as the handler to a POST request to create a new healthcheck.
func CreateHealthcheck(c *gin.Context) {
	// swagger:route POST /v1/healthcheck healthchecks createHealthcheck
	//
	// Creates a new healthcheck to monitor.
	//
	//     Consumes:
	//     - application/json
	//     - application/yaml
	//
	//     Produces:
	//     - application/json
	//     - application/yaml
	//
	//     Schemes: https
	//
	//     Responses:
	//       201: createHealthcheckResponse
	createHealthcheckParams := CreateHealthcheckParams{}
	if err := c.Bind(&createHealthcheckParams); err != nil {
		log.Error(err)
		return
	}

	healthcheck, err := healthcheck.CreateOrUpdateHealthcheck(createHealthcheckParams.Healthcheck)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	createHealthcheckResponse := CreateHealthcheckResponse{
		Healthcheck: healthcheck,
	}

	c.JSON(http.StatusCreated, createHealthcheckResponse)
}
