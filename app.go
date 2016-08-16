package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/premkit/healthcheck/commands"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.Execute()
}
