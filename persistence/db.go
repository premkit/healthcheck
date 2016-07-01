package persistence

import (
	"os"
	"path"

	"github.com/premkit/healthcheck/log"

	"github.com/boltdb/bolt"
)

// DB is the lazy-loaded reference to the BoltDB instance.  Use the GetDB() function to obtain this.
var DB *bolt.DB

// GetDB returns the singleton instance of the BoltDB connection.  This is not a threadsafe object,
// but transactions are.  Any caller using this object should use a transaction.
func GetDB() (*bolt.DB, error) {
	if DB != nil {
		return DB, nil
	}

	// Use the environment variable for the datadir, if set
	dataDirectory := os.Getenv("DATA_DIRECTORY")
	if dataDirectory == "" {
		dataDirectory = "/data"
	}

	if err := os.MkdirAll(dataDirectory, 0600); err != nil {
		log.Error(err)
		return nil, err
	}

	path := path.Join(dataDirectory, "healthcheck.db")
	conn, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err := initializeDatabase(conn); err != nil {
		return nil, err
	}

	DB = conn
	return DB, nil
}

func initializeDatabase(conn *bolt.DB) error {
	// Perform some initialization
	err := conn.Update(func(tx *bolt.Tx) error {
		// Create the default buckets
		_, err := tx.CreateBucketIfNotExists([]byte("Index"))
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})

	return err
}
