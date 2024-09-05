package store

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/core"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/oklog/ulid/v2"
	bolt "go.etcd.io/bbolt"
)

type VersionStore interface {
	CreateVersion(ctx context.Context, pages map[string]ulid.ULID, theme ulid.ULID) (core.Version, error)
	GetVersion(ctx context.Context, uid ulid.ULID) (core.Version, error)
	SetPage(ctx context.Context, uid ulid.ULID, previousPath string, path string, pageUID ulid.ULID) (core.Version, error)
	Clone(ctx context.Context, uid ulid.ULID) (core.Version, error)
	GetPageVersions(ctx context.Context, pageUID ulid.ULID) (mapset.Set[ulid.ULID], error)
}

type versionStore struct {
	db           documentDB[core.Version]
	pageVersions PageVersionIndex
}

func (s versionStore) CreateVersion(ctx context.Context, pages map[string]ulid.ULID, theme ulid.ULID) (core.Version, error) {
	version := core.Version{
		UID:   ulid.Make(),
		Pages: pages,
		Theme: theme,
	}
	if err := s.db.Save(ctx, bucketVersions, version.UID.String(), version); err != nil {
		return core.Version{}, err
	}
	if err := s.pageVersions.CreateVersion(ctx, &version); err != nil {
		return core.Version{}, err
	}
	return version, nil
}

func (s versionStore) GetVersion(ctx context.Context, uid ulid.ULID) (core.Version, error) {
	return s.db.One(ctx, bucketVersions, uid.String())
}

func (s versionStore) SetPage(ctx context.Context, uid ulid.ULID, previousPath string, path string, pageUID ulid.ULID) (core.Version, error) {
	version, err := s.db.One(ctx, bucketVersions, uid.String())
	if err != nil {
		return core.Version{}, err
	}
	if previousPath != "" {
		delete(version.Pages, previousPath)
	}
	version.Pages[path] = pageUID
	if err := s.db.Save(ctx, bucketVersions, uid.String(), version); err != nil {
		return core.Version{}, err
	}
	if err := s.pageVersions.Add(ctx, pageUID, uid); err != nil {
		return core.Version{}, nil
	}
	return version, nil
}

func (s versionStore) Clone(ctx context.Context, uid ulid.ULID) (core.Version, error) {
	version, err := s.db.One(ctx, bucketVersions, uid.String())
	if err != nil {
		return core.Version{}, err
	}
	version.UID = ulid.Make()
	if err := s.db.Save(ctx, bucketVersions, version.UID.String(), version); err != nil {
		return core.Version{}, err
	}
	return version, nil
}

func (s versionStore) GetPageVersions(ctx context.Context, pageUID ulid.ULID) (mapset.Set[ulid.ULID], error) {
	return s.pageVersions.GetVersions(ctx, pageUID)
}

func NewVersionStore(db *bolt.DB, pageVersions PageVersionIndex) VersionStore {
	return &versionStore{
		db:           docDB[core.Version]{db},
		pageVersions: pageVersions,
	}
}
