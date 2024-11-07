package store

import (
	"bytes"
	"context"
	"encoding/gob"
	"strings"

	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/oklog/ulid/v2"
	bolt "go.etcd.io/bbolt"
)

type documentDB[T any] interface {
	One(ctx context.Context, bucket string, key string) (T, error)
	Many(ctx context.Context, bucket string, start *string, count int) ([]T, *string, error)
	Save(ctx context.Context, bucket string, key string, item T) error
}

type docDB[T any] struct {
	db *bolt.DB
}

func getStringKey(input *ulid.ULID) *string {
	if input == nil {
		return nil
	}
	str := input.String()
	return &str
}

func getULIDKey(input *string) (*ulid.ULID, error) {
	if input == nil {
		return nil, nil
	}
	id, err := ulid.Parse(*input)
	if err != nil {
		return nil, err
	}
	return &id, nil
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
		decoder := gob.NewDecoder(bytes.NewBuffer(raw))
		if err := decoder.Decode(&item); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return item, err
	}
	return item, nil
}

func (d docDB[T]) Many(ctx context.Context, bucket string, start *string, count int) ([]T, *string, error) {
	items := make([]T, 0)
	var nextItemKey strings.Builder
	if err := transactCtx(ctx, d.db, false, func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return core.ErrBucketNotFound.New("bucket %s does not exist", bucket)
		}
		c := b.Cursor()
		var k, v []byte
		if start == nil {
			k, v = c.Last()
		} else {
			k, v = c.Seek([]byte(*start))
		}
		for v != nil && len(items) < count {
			var item T
			decoder := gob.NewDecoder(bytes.NewBuffer(v))
			if err := decoder.Decode(&item); err != nil {
				return err
			}
			items = append(items, item)
			k, v = c.Prev()
		}
		if k != nil {
			nextItemKey.Write(k)
		}
		return nil
	}); err != nil {
		return nil, nil, err
	}
	var nextItemKeyString *string
	if nextItemKey.Len() > 0 {
		collected := nextItemKey.String()
		nextItemKeyString = &collected
	}
	return items, nextItemKeyString, nil
}

func (d docDB[T]) Save(ctx context.Context, bucket string, key string, item T) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(&item); err != nil {
		return err
	}
	return transactCtx(ctx, d.db, true, func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return core.ErrBucketNotFound.New("bucket %s does not exist", bucket)
		}
		return b.Put([]byte(key), buffer.Bytes())
	})
}
