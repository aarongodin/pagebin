package store

import (
	"context"
	"encoding/json"

	"github.com/aarongodin/pagebin/pkg/core"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/oklog/ulid/v2"
	bolt "go.etcd.io/bbolt"
)

type PageVersionIndex interface {
	GetVersions(ctx context.Context, pageUID ulid.ULID) (mapset.Set[ulid.ULID], error)
	CreateVersion(ctx context.Context, version *core.Version) error
	Add(ctx context.Context, pageUID ulid.ULID, versionUID ulid.ULID) error
	Remove(ctx context.Context, pageUID ulid.ULID, versionUID ulid.ULID) error
}

type pageVersionIndex struct {
	db *bolt.DB
}

func (i pageVersionIndex) GetVersions(ctx context.Context, pageUID ulid.ULID) (mapset.Set[ulid.ULID], error) {
	versions := mapset.NewSet[ulid.ULID]()
	if err := transactCtx(ctx, i.db, false, func(tx *bolt.Tx) error {
		b, err := getIndexBucket(tx, bucketIndexPageVersions)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(b.Get(pageUID.Bytes()), &versions); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return versions, nil
}

func (i pageVersionIndex) Add(ctx context.Context, pageUID ulid.ULID, versionUID ulid.ULID) error {
	return i.modify(ctx, pageUID, func(versions mapset.Set[ulid.ULID]) {
		versions.Add(versionUID)
	})
}

func (i pageVersionIndex) Remove(ctx context.Context, pageUID ulid.ULID, versionUID ulid.ULID) error {
	return i.modify(ctx, pageUID, func(versions mapset.Set[ulid.ULID]) {
		versions.Remove(versionUID)
	})
}

func (i pageVersionIndex) modify(ctx context.Context, pageUID ulid.ULID, fn func(versions mapset.Set[ulid.ULID])) error {
	return transactCtx(ctx, i.db, true, func(tx *bolt.Tx) error {
		b, err := getIndexBucket(tx, bucketIndexPageVersions)
		if err != nil {
			return err
		}
		versions := mapset.NewSet[ulid.ULID]()
		stored := b.Get(pageUID.Bytes())
		if stored != nil {
			if err := json.Unmarshal(stored, &versions); err != nil {
				return err
			}
		}
		fn(versions)
		raw, err := json.Marshal(&versions)
		if err != nil {
			return err
		}
		return b.Put(pageUID.Bytes(), raw)
	})
}

func (i pageVersionIndex) CreateVersion(ctx context.Context, version *core.Version) error {
	return transactCtx(ctx, i.db, true, func(tx *bolt.Tx) error {
		b, err := getIndexBucket(tx, bucketIndexPageVersions)
		if err != nil {
			return err
		}
		for _, pageUID := range version.Pages {
			versions := mapset.NewSet[ulid.ULID]()
			stored := b.Get(pageUID.Bytes())
			if stored != nil {
				if err := json.Unmarshal(stored, &versions); err != nil {
					return err
				}
			}
			versions.Add(version.UID)
			raw, err := json.Marshal(&versions)
			if err != nil {
				return err
			}
			if err := b.Put(pageUID.Bytes(), raw); err != nil {
				return err
			}
		}
		return nil
	})
}

func NewPageVersionIndex(db *bolt.DB) PageVersionIndex {
	return &pageVersionIndex{db}
}
