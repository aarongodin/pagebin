package store

import (
	"github.com/aarongodin/pagebin/pkg/core"
	bolt "go.etcd.io/bbolt"
)

func getIndexBucket(tx *bolt.Tx, name string) (*bolt.Bucket, error) {
	i := tx.Bucket([]byte(bucketIndex))
	if i == nil {
		return nil, core.ErrBucketNotFound.New("bucket %s does not exist", bucketIndex)
	}
	b := i.Bucket([]byte(name))
	if b == nil {
		return nil, core.ErrBucketNotFound.New("bucket %s does not exist", name)
	}
	return b, nil
}
