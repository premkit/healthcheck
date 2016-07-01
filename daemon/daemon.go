package daemon

import (
	"github.com/gin-gonic/gin"
)

// Run is the main entrypoint of this daemon.
func Run() {
	go func() {
		r := gin.Default()

		r.RedirectTrailingSlash = false

		apiV1Group := r.Group("/v1")

		apiV1Group.POST("/healthcheck", CreateHealthcheck)
		apiV1Group.GET("/healthchecks", ListHealthchecks)

		r.Run(":80")
	}()
}
