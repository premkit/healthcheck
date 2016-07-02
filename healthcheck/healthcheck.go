package healthcheck

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/premkit/healthcheck/httpcheck"
	"github.com/premkit/healthcheck/log"
	"github.com/premkit/healthcheck/persistence"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-multierror"
)

/*
  - name: Website
    checks:
      - type: http
        spec:
          method: "GET"
          options: ['no_follow_redirects', 'ignore_insecure']
          status_codes_available: [200]
          uri: "http://$CONTAINER_PYTHONWEB_MARCREFERENCE_PRIVATEIP:$CONTAINER_PYTHONWEB_MARCREFERENCE_HOSTPORT_5000/_ping"
*/

// Healthcheck represents one entire healthcheck.  A healthcheck has a name and contains
// one or more steps that make up the healthcheck.
type Healthcheck struct {
	Name       string                       `json:"name"`
	HTTPChecks []*httpcheck.HTTPHealthcheck `json:"http_checks"`
}

func validateHealthcheck(healthcheck *Healthcheck) error {
	var result *multierror.Error

	if healthcheck == nil {
		err := errors.New("Healthcheck cannot be nil")
		log.Error(err)
		result = multierror.Append(result, err)
		return result.ErrorOrNil()
	}

	if healthcheck.Name == "" {
		err := errors.New("Name is required")
		log.Error(err)
		result = multierror.Append(result, err)
	}

	if len(healthcheck.HTTPChecks) == 0 {
		err := errors.New("At least one step is required")
		log.Error(err)
		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
}

// CreateOrUpdateHealthcheck will create or update the specified healthcheck in the database.
func CreateOrUpdateHealthcheck(healthcheck *Healthcheck) (*Healthcheck, error) {
	if err := validateHealthcheck(healthcheck); err != nil {
		return nil, err
	}

	existing, err := getHealthcheckByName([]byte(healthcheck.Name))
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return updateHealthcheck(existing, healthcheck)
	}

	return createHealthcheck(healthcheck)
}

func updateHealthcheck(existing *Healthcheck, updated *Healthcheck) (*Healthcheck, error) {
	db, err := persistence.GetDB()
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		log.Debugf("Updating existing healthcheck %s", existing.Name)

		// Get the bucket
		bucket := tx.Bucket([]byte(existing.Name))

		if bucket == nil {
			err := errors.New("Bucket not found")
			log.Error(err)
			return err
		}
		// Update the properties
		if err := bucket.Put([]byte("name"), []byte(updated.Name)); err != nil {
			log.Error(err)
			return err
		}

		// Update httpchecks
		cursor := bucket.Cursor()
		prefix := []byte("httpcheck.")
		for k, _ := cursor.Seek(prefix); bytes.HasPrefix(k, prefix); k, _ = cursor.Next() {
			if err := bucket.Delete(k); err != nil {
				log.Error(err)
				return err
			}
		}

		for i, httpCheck := range updated.HTTPChecks {
			b, err := httpCheck.Serialize()
			if err != nil {
				return err
			}

			if err := bucket.Put([]byte(fmt.Sprintf("httpcheck.%d", i)), b); err != nil {
				log.Error(err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return updated, nil
}

func createHealthcheck(healthcheck *Healthcheck) (*Healthcheck, error) {
	db, err := persistence.GetDB()
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		log.Debugf("Creating bucket %q for new healthcheck", healthcheck.Name)

		indexBucket := tx.Bucket([]byte("Index"))
		if err := indexBucket.Put([]byte(healthcheck.Name), []byte(healthcheck.Name)); err != nil {
			log.Error(err)
			return err
		}

		// Create the new bucket to hold this healthcheck
		bucket, err := tx.CreateBucket([]byte(healthcheck.Name))
		if err != nil {
			log.Error(err)
			return err
		}

		// Serialize this healthcheck to the bucket
		if err := bucket.Put([]byte("name"), []byte(healthcheck.Name)); err != nil {
			log.Error(err)
			return err
		}

		for i, httpCheck := range healthcheck.HTTPChecks {
			b, err := httpCheck.Serialize()
			if err != nil {
				return err
			}

			if err := bucket.Put([]byte(fmt.Sprintf("httpcheck.%d", i)), b); err != nil {
				log.Error(err)
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return healthcheck, nil
}

// ListHealthchecks will return an array of all known healthchecks, read from the database.
func ListHealthchecks() ([]*Healthcheck, error) {
	db, err := persistence.GetDB()
	if err != nil {
		return nil, err
	}

	healthchecks := make([]*Healthcheck, 0, 0)

	err = db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			if string(name) == "Index" {
				return nil
			}

			healthcheck, err := getHealthcheckByName(name)
			if err != nil {
				return err
			}

			healthchecks = append(healthchecks, healthcheck)

			return nil
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return healthchecks, nil
}

func getHealthcheckByName(name []byte) (*Healthcheck, error) {
	db, err := persistence.GetDB()
	if err != nil {
		return nil, err
	}

	var healthcheck *Healthcheck

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(name)
		if bucket == nil {
			return nil
		}

		healthcheck = &Healthcheck{
			Name: string(bucket.Get([]byte("name"))),
		}

		// Get the individual checks
		cursor := bucket.Cursor()
		prefix := []byte("httpcheck.")
		for k, v := cursor.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			httpCheck, err := httpcheck.Deserialize(v)
			if err != nil {
				return err
			}

			if healthcheck.HTTPChecks == nil {
				healthcheck.HTTPChecks = make([]*httpcheck.HTTPHealthcheck, 0, 0)
			}

			healthcheck.HTTPChecks = append(healthcheck.HTTPChecks, httpCheck)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return healthcheck, nil
}
