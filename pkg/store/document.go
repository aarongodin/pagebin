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

func transactContext(db *bolt.DB, writable bool, ctx context.Context, fn func(tx *bolt.Tx) error) error {
	tx, err := db.Begin(writable)
	if err != nil {
		return err
	}
	done := make(chan error, 1)

	go func() {
		done <- fn(tx)
	}()

	select {
	case <-ctx.Done():
		return tx.Rollback()
	case err := <-done:
		if writable {
			if err == nil {
				return tx.Commit()
			} else {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return err
			}
		} else {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
}

func (d docDB[T]) One(ctx context.Context, bucket string, key string) (T, error) {
	var item T
	if err := transactContext(d.db, false, ctx, func(tx *bolt.Tx) error {
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
	return transactContext(d.db, true, ctx, func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return core.ErrBucketNotFound.New("bucket %s does not exist", bucket)
		}
		return b.Put([]byte(key), raw)
	})
}
