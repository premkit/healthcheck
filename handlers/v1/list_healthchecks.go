package v1

import (
	"net/http"

	"github.com/premkit/healthcheck/healthcheck"

	"github.com/gin-gonic/gin"
)

// ListHealthchecksResponse represents the response to a listHealthchecks call. This response
// includes a list of all available healthchecks.
// swagger:response listHealthchecksResponse
type ListHealthchecksResponse struct {
	// Healthcheck
	// In: body
	Healthchecks []*healthcheck.Healthcheck `json:"healthchecks"`
}

// ListHealthchecks is called as the handler to a GET request to list healthchecks.
func ListHealthchecks(c *gin.Context) {
	// swagger:route GET /v1/healthchecks healthchecks listHealthchecks
	//
	// Lists all available healthchecks
	//
	//     Produces:
	//     - application/json
	//     - application/yaml
	//
	//     Schemes: https
	//
	//     Responses:
	//       200: listHealthchecksResponse

	healthchecks, err := healthcheck.ListHealthchecks()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	listHealthchecksResponse := ListHealthchecksResponse{
		Healthchecks: healthchecks,
	}

	c.JSON(http.StatusOK, listHealthchecksResponse)
}
