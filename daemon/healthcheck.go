package daemon

import (
	"net/http"

	"github.com/premkit/healthcheck/healthcheck"
	"github.com/premkit/healthcheck/schema"

	"github.com/gin-gonic/gin"
)

func CreateHealthcheck(c *gin.Context) {
	request := schema.CreateHealthcheckRequest{}
	if err := c.Bind(&request); err != nil {
		// TODO log
		return
	}

	healthcheck, err := healthcheck.CreateHealthcheck(request.Healthcheck)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := schema.CreateHealthcheckResponse{
		Healthcheck: healthcheck,
	}

	c.JSON(http.StatusOK, response)
}
