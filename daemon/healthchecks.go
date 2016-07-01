package daemon

import (
	"net/http"

	"github.com/premkit/healthcheck/healthcheck"
	"github.com/premkit/healthcheck/schema"

	"github.com/gin-gonic/gin"
)

func ListHealthchecks(c *gin.Context) {
	healthchecks, err := healthcheck.ListHealthchecks()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := schema.ListHealthchecksResponse{
		Healthchecks: healthchecks,
	}

	c.JSON(http.StatusOK, response)
}
