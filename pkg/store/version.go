package store

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/oklog/ulid/v2"
	bolt "go.etcd.io/bbolt"
)

type VersionStore interface {
	CreateVersion(ctx context.Context, pages []ulid.ULID, theme ulid.ULID) (core.Version, error)
	GetVersion(ctx context.Context, uid ulid.ULID) (core.Version, error)
}

type versionStore struct {
	db documentDB[core.Version]
}

func (s versionStore) CreateVersion(ctx context.Context, pages []ulid.ULID, theme ulid.ULID) (core.Version, error) {
	version := core.Version{
		UID:   ulid.Make(),
		Pages: pages,
		Theme: theme,
	}
	if err := s.db.Save(ctx, bucketVersions, version.UID.String(), version); err != nil {
		return core.Version{}, err
	}
	return version, nil
}

func (s versionStore) GetVersion(ctx context.Context, uid ulid.ULID) (core.Version, error) {
	return s.db.One(ctx, bucketVersions, uid.String())
}

func NewVersionStore(db *bolt.DB) VersionStore {
	return &versionStore{db: docDB[core.Version]{db}}
}
