package server

import (
	"fmt"

	"github.com/premkit/healthcheck/handlers/v1"
	"github.com/premkit/healthcheck/log"

	"github.com/gin-gonic/gin"
)

// Run is the main entrypoint of this daemon.
func Run(config *Config) {
	go func() {
		r := gin.Default()

		r.RedirectTrailingSlash = false

		apiV1Group := r.Group("/v1")

		apiV1Group.GET("/healthchecks", v1.ListHealthchecks)
		apiV1Group.POST("/healthcheck", v1.CreateHealthcheck)

		log.Infof("Listening on port %d for http connections", config.HTTPPort)
		r.Run(fmt.Sprintf(":%d", config.HTTPPort))
	}()
}
