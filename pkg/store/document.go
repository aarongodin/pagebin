package store

import (
	"context"
	"encoding/json"

	"github.com/aarongodin/pagebin/pkg/core"
	bolt "go.etcd.io/bbolt"
)

type documentDB[T any] interface {
	One(ctx context.Context, bucket string, key string) (T, error)
	Save(ctx context.Context, bucket string, key string, item T) error
}

type docDB[T any] struct {
	db *bolt.DB
}

func (d docDB[T]) One(ctx context.Context, bucket string, key string) (T, error) {
	var item T
	if err := transactCtx(ctx, d.db, false, func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return core.ErrBucketNotFound.New("bucket %s does not exist", bucket)
		}
		raw := b.Get([]byte(key))
		if raw == nil {
			return core.ErrItemNotFound.New("item %s/%s not found", bucket, key)
		}
		if err := json.Unmarshal(raw, &item); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return item, err
	}
	return item, nil
}

func (d docDB[T]) Save(ctx context.Context, bucket string, key string, item T) error {
	raw, err := json.Marshal(&item)
	if err != nil {
		return err
	}
	return transactCtx(ctx, d.db, true, func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return core.ErrBucketNotFound.New("bucket %s does not exist", bucket)
		}
		return b.Put([]byte(key), raw)
	})
}
