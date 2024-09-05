package store

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/oklog/ulid/v2"
	bolt "go.etcd.io/bbolt"
)

var (
	keySite = "site"
)

type SiteStore interface {
	GetSite(ctx context.Context) (core.Site, error)
	CreateSite(ctx context.Context, title string, version ulid.ULID, nextVersion ulid.ULID) (core.Site, error)
	UpdateSite(ctx context.Context, title string) (core.Site, error)
}

type siteStore struct {
	db documentDB[core.Site]
}

func (s siteStore) GetSite(ctx context.Context) (core.Site, error) {
	return s.db.One(ctx, bucketApp, keySite)
}

func (s siteStore) CreateSite(ctx context.Context, title string, version ulid.ULID, nextVersion ulid.ULID) (core.Site, error) {
	site := core.Site{
		UID:         ulid.Make(),
		Title:       title,
		Version:     version,
		NextVersion: nextVersion,
	}
	if err := s.db.Save(ctx, bucketApp, keySite, site); err != nil {
		return core.Site{}, err
	}
	return site, nil
}

func (s siteStore) UpdateSite(ctx context.Context, title string) (core.Site, error) {
	site, err := s.db.One(ctx, bucketApp, keySite)
	if err != nil {
		return core.Site{}, err
	}
	if title != "" {
		site.Title = title
	}
	if err := s.db.Save(ctx, bucketApp, keySite, site); err != nil {
		return core.Site{}, err
	}
	return site, nil
}

func NewSiteStore(db *bolt.DB) SiteStore {
	return &siteStore{db: docDB[core.Site]{db}}
}
