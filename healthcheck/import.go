package healthcheck

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/premkit/healthcheck/httpcheck"
	"github.com/premkit/healthcheck/log"

	"gopkg.in/yaml.v2"
)

// ImportServiceFile will parse the file at the given location and load any service
// definitions found into the running system
// Note that this will merge any existing definitions that match the same ids, which
// will allow history to be retained.
func ImportServiceFile(contentType string, filepath string) ([]*Healthcheck, error) {
	healthchecks, err := readServiceFile(contentType, filepath)
	if err != nil {
		return nil, err
	}

	for i, h := range healthchecks {
		healthcheck, err := CreateOrUpdateHealthcheck(h)
		if err != nil {
			return nil, err
		}

		healthchecks[i] = healthcheck
	}

	return healthchecks, nil
}

func readServiceFile(contentType string, filepath string) ([]*Healthcheck, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	switch contentType {
	case "application/json":
		return unmarshalJSONService(data)

	case "application/yaml":
		return unmarshalYAMLService(data)
	}

	err = fmt.Errorf("Unknown content type: %q", contentType)
	log.Error(err)
	return nil, err
}

func unmarshalJSONService(data []byte) ([]*Healthcheck, error) {
	var healthchecks []*Healthcheck

	var unmarshal []*healthcheckUnmarshal
	if err := json.Unmarshal(data, &unmarshal); err != nil {
		log.Error(err)
		return nil, err
	}

	for _, h := range unmarshal {
		healthcheck := Healthcheck{
			Name: h.Name,
		}

		if err := unmarshalChecks(&healthcheck, h.Checks); err != nil {
			return nil, err
		}

		healthchecks = append(healthchecks, &healthcheck)
	}

	return healthchecks, nil
}

func unmarshalYAMLService(data []byte) ([]*Healthcheck, error) {
	var healthchecks []*Healthcheck

	var unmarshal []*healthcheckUnmarshal
	if err := yaml.Unmarshal(data, &unmarshal); err != nil {
		log.Infof("aaa")
		log.Error(err)
		return nil, err
	}

	for _, h := range unmarshal {
		healthcheck := Healthcheck{
			Name: h.Name,
		}

		if err := unmarshalChecks(&healthcheck, h.Checks); err != nil {
			log.Infof("bbb")
			return nil, err
		}

		healthchecks = append(healthchecks, &healthcheck)
	}

	return healthchecks, nil
}

func unmarshalChecks(healthcheck *Healthcheck, checks []*checkUnmarshal) error {
	for _, check := range checks {
		switch check.Type {
		case "http":
			httpCheck, err := unmarshalHTTPCheck(check.Spec)
			if err != nil {
				return err
			}

			healthcheck.HTTPChecks = append(healthcheck.HTTPChecks, httpCheck)
		default:
			err := fmt.Errorf("Unknown check type: %q", check.Type)
			log.Error(err)
			return err
		}
	}

	return nil
}

func unmarshalHTTPCheck(check map[string]interface{}) (*httpcheck.HTTPHealthcheck, error) {
	b, err := json.Marshal(check)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	result := httpcheck.HTTPHealthcheck{}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Error(err)
		return nil, err
	}

	return &result, nil
}
