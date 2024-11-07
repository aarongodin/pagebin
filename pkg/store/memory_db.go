package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

func withTestDB(t *testing.T, bucketName string, cb func(db *bolt.DB)) {
	file, err := os.CreateTemp("", "pagebin-testdb-")
	require.NoError(t, err)
	defer os.Remove(file.Name())
	db, err := bolt.Open(file.Name(), 0600, nil)
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucketName))
		return err
	}))
	cb(db)
}
