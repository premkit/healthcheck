package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/premkit/healthcheck/log"

	"github.com/parnurzeal/gorequest"
)

// This package interfaces with the upstream premkit router

// Register will register this server with the upstream router.
func Register() error {
	if os.Getenv("PREMKIT_ROUTER") == "" {
		log.Infof("PREMKIT_ROUTER not defined. Not going to register with an upstream.")
		return nil
	}

	// TODO once the code is in github, use a reference to the schema instead of hard coding this json object
	body := fmt.Sprintf(`
{
	"service": {
		"name": "healthcheck",
		"path": "healthcheck",
		"upstreams": [
			%q
		]
	}
}`, os.Getenv("ADVERTISE_ADDRESS"))

	uri := fmt.Sprintf("%s/premkit/v1/service", os.Getenv("PREMKIT_ROUTER"))
	request := gorequest.New()
	resp, _, errs := request.Post(uri).
		Send(body).
		End()

	if len(errs) > 0 {
		err := fmt.Errorf("Error(s) registering with router: %#v", errs)
		log.Error(err)
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		err := fmt.Errorf("Unexpected response from router: %d", resp.StatusCode)
		log.Error(err)
		return err
	}

	return nil
}
