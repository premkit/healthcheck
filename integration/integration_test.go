package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/premkit/healthcheck/handlers/v1"
	"github.com/premkit/healthcheck/healthcheck"
	"github.com/premkit/healthcheck/httpcheck"
	"github.com/premkit/healthcheck/schema"
	"github.com/premkit/healthcheck/server"

	"github.com/parnurzeal/gorequest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	baseURL        string
	dbPathOriginal string
)

func TestMain(m *testing.M) {
	// Initialize a temp boltdb folder
	err := setup()
	if err != nil {
		os.Exit(1)
	}
	defer teardown()

	// Run the tests
	result := m.Run()

	os.Exit(result)
}

func setup() error {
	config := &server.Config{
		HTTPPort: 8881, // This has to be high enough for Circle to allow it
	}
	server.Run(config)

	dbPathOriginal = os.Getenv("DATA_DIRECTORY")
	dirName, err := ioutil.TempDir("", "integration")
	if err != nil {
		return err
	}
	os.Setenv("DATA_DIRECTORY", dirName)

	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = "8881"
	}
	baseURL = fmt.Sprintf("http://localhost:%s", port)

	return nil
}

func teardown() {
	os.Setenv("DATA_DIRECTORY", dbPathOriginal)
}

// Create a functional healthcheck, list them, and run all.
func TestHealthchecks(t *testing.T) {
	// Create a healthcheck
	isGoogleOK := httpcheck.HTTPHealthcheck{
		URI:                  "https://www.google.com",
		Method:               "GET",
		StatusCodesAvailable: []int{http.StatusOK},
	}

	healthcheck := healthcheck.Healthcheck{
		Name:       "INTEGRATION",
		HTTPChecks: []*httpcheck.HTTPHealthcheck{&isGoogleOK},
	}

	createRequest := v1.CreateHealthcheckParams{
		Healthcheck: &healthcheck,
	}

	request := gorequest.New()
	resp, body, errs := request.Post(fmt.Sprintf("%s/v1/healthcheck", baseURL)).
		Send(createRequest).
		End()

	assert.Nil(t, errs, "there should be no errors")
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	createResponse := v1.CreateHealthcheckResponse{}
	err := json.Unmarshal([]byte(body), &createResponse)
	require.NoError(t, err)

	assert.Equal(t, createResponse.Healthcheck.Name, "INTEGRATION", "name should be 'INTEGRATION'")

	// List healthchecks
	resp, body, errs = gorequest.New().Get(fmt.Sprintf("%s/v1/healthchecks", baseURL)).End()

	assert.Equal(t, 0, len(errs), "there should be no errors")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	listResponse := schema.ListHealthchecksResponse{}
	err = json.Unmarshal([]byte(body), &listResponse)
	require.NoError(t, err)

	assert.Equal(t, 1, len(listResponse.Healthchecks), "there should be 1 healthcheck")
	assert.Equal(t, "INTEGRATION", listResponse.Healthchecks[0].Name, "the healthcheck name should be 'INTEGRATION'")

}
