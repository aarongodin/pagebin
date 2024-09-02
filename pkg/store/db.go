package store

import (
	"github.com/aarongodin/pagebin/pkg/config"
	bolt "go.etcd.io/bbolt"
)

var (
	bucketApp      = "app"
	bucketThemes   = "themes"
	bucketPages    = "pages"
	bucketVersions = "versions"
	buckets        = [...]string{bucketApp, bucketThemes, bucketPages, bucketVersions}
)

func NewDB(rc *config.RuntimeConfig) (*bolt.DB, error) {
	db, err := bolt.Open(rc.DatabaseFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		for _, b := range buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(b)); err != nil {
				return err
			}
		}
		return nil
	})

	return db, nil
}
