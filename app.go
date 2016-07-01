package main

import (
	"math/rand"
	"time"

	"github.com/premkit/healthcheck/daemon"
	"github.com/premkit/healthcheck/router"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	// Import healthchecks from the bootstrap config.

	// Start the HTTP server.
	daemon.Run()

	// Register with the reverse proxy, if the reverse proxy is available
	router.Register()

	<-make(chan int)
}
