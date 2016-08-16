package healthcheck

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: these tests need to run without overwriting the primary data store

func TestReadYAMLFile(t *testing.T) {
	dirName, err := ioutil.TempDir("", "healthcheck-test")
	require.NoError(t, err)
	defer os.RemoveAll(dirName)

	filename := path.Join(dirName, "service.yaml")

	err = ioutil.WriteFile(filename, []byte(simpleServiceYAML), 0444)
	require.NoError(t, err)

	healthchecks, err := readServiceFile("application/yaml", filename)
	require.NoError(t, err)
	require.NotNil(t, healthchecks)
	require.Equal(t, 1, len(healthchecks))

	service := healthchecks[0]
	assert.Equal(t, "service1", service.Name)

	assert.Equal(t, 1, len(service.HTTPChecks))

	check := service.HTTPChecks[0]
	assert.Equal(t, "GET", check.Method)
	assert.Equal(t, "http://localhost", check.URI)
	assert.Equal(t, 2, len(check.Options))
	assert.Equal(t, 1, len(check.StatusCodesAvailable))
	assert.Equal(t, 200, check.StatusCodesAvailable[0])
	assert.Nil(t, check.StatusCodesDegraded)
	assert.Nil(t, check.StatusCodesUnavailable)
}

func TestReadJSONFile(t *testing.T) {
	dirName, err := ioutil.TempDir("", "healthcheck-test")
	require.NoError(t, err)
	defer os.RemoveAll(dirName)

	filename := path.Join(dirName, "service.json")

	err = ioutil.WriteFile(filename, []byte(simpleServiceJSON), 0444)
	require.NoError(t, err)

	healthchecks, err := readServiceFile("application/json", filename)
	require.NoError(t, err)
	require.NotNil(t, healthchecks)
	require.Equal(t, 1, len(healthchecks))

	service := healthchecks[0]
	assert.Equal(t, "service1", service.Name)

	assert.Equal(t, 1, len(service.HTTPChecks))

	check := service.HTTPChecks[0]
	assert.Equal(t, "POST", check.Method)
	assert.Equal(t, "http://localhost", check.URI)
	assert.Equal(t, 2, len(check.Options))
	assert.Equal(t, 1, len(check.StatusCodesAvailable))
	assert.Equal(t, 200, check.StatusCodesAvailable[0])
	assert.Nil(t, check.StatusCodesDegraded)
	assert.Nil(t, check.StatusCodesUnavailable)
}

func TestImportService(t *testing.T) {
	dirName, err := ioutil.TempDir("", "healthcheck-test")
	require.NoError(t, err)
	defer os.RemoveAll(dirName)

	filename := path.Join(dirName, "service.yaml")
	err = ioutil.WriteFile(filename, []byte(simpleServiceYAML), 0444)
	require.NoError(t, err)

	imported, err := ImportServiceFile("application/yaml", filename)
	require.NoError(t, err)
	assert.Equal(t, 1, len(imported))
	listed, err := ListHealthchecks()
	require.NoError(t, err)
	assert.Equal(t, len(imported), len(listed))
	assert.Equal(t, "service1", listed[0].Name)
	require.Equal(t, 1, len(listed[0].HTTPChecks))
	assert.Equal(t, "GET", listed[0].HTTPChecks[0].Method)

	filename = path.Join(dirName, "service.json")
	err = ioutil.WriteFile(filename, []byte(simpleServiceJSON), 0444)
	require.NoError(t, err)

	imported2, err := ImportServiceFile("application/json", filename)
	require.NoError(t, err)
	assert.Equal(t, 1, len(imported2))
	listed2, err := ListHealthchecks()
	require.NoError(t, err)
	assert.Equal(t, 1, len(listed2))
	assert.Equal(t, "service1", listed2[0].Name)
	require.Equal(t, 1, len(listed2[0].HTTPChecks))
	assert.Equal(t, "POST", listed2[0].HTTPChecks[0].Method)
}
