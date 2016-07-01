package healthcheck

import (
	"errors"

	"github.com/premkit/healthcheck/log"
	"github.com/premkit/healthcheck/persistence"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-multierror"
	"github.com/tuvistavie/securerandom"
)

type HealthcheckResponse int

const (
	HealthcheckResponseUnknown HealthcheckResponse = iota
	HealthcheckResponseAvailable
	HealthcheckResponseDegraded
	HealthcheckResponseUnavailable
)

type Healthcheck struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	Steps []HealthcheckStep `json:"steps"`
}

type HealthcheckStep interface {
	RunSynchronously() HealthcheckStepResult
}

type HealthcheckStepResult interface {
	Status() HealthcheckResponse
}

func validateHealthcheck(healthcheck *Healthcheck) error {
	var result *multierror.Error

	if healthcheck.Name == "" {
		result = multierror.Append(result, errors.New("Name is required"))
	}

	if len(healthcheck.Steps) == 0 {
		result = multierror.Append(result, errors.New("At least one step is required"))
	}

	return result.ErrorOrNil()
}

// CreateHealthcheck writes the new healthcheck to the database and initializes it to be available.  After calling this,
// the returning healthcheck object will have an ID and be immediately available.
func CreateHealthcheck(healthcheck *Healthcheck) (*Healthcheck, error) {
	if err := validateHealthcheck(healthcheck); err != nil {
		return nil, err
	}

	id, err := securerandom.Hex(32)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	healthcheck.ID = id

	db, err := persistence.GetDB()
	if err != nil {
		return nil, err
	}

	err = db.Batch(func(tx *bolt.Tx) error {
		// Create the new bucket to hold this healthcheck
		bucket, err := tx.CreateBucket([]byte(id))
		if err != nil {
			log.Error(err)
			return err
		}

		// Serialize this healthcheck to the bucket
		if err := bucket.Put([]byte("name"), []byte(healthcheck.Name)); err != nil {
			log.Error(err)
			return err
		}
		indexBucket := tx.Bucket([]byte("Index"))
		if err := indexBucket.Put([]byte(id), []byte(id)); err != nil {
			log.Error(err)
			return err
		}

		return nil
	})

	return healthcheck, nil
}

func ListHealthchecks() ([]*Healthcheck, error) {
	return nil, nil
}
