package store

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/config"
	bolt "go.etcd.io/bbolt"
)

type Store interface {
	DB() *bolt.DB
	StartTx(ctx context.Context, writable bool) (context.Context, error)
	EndTx(ctx context.Context, txErr error) error
	Close(ctx context.Context) error
	Sites() SiteStore
	Pages() PageStore
	Versions() VersionStore
	Themes() ThemeStore
	Blobs() BlobStore
}

type store struct {
	db       *bolt.DB
	sites    SiteStore
	pages    PageStore
	versions VersionStore
	themes   ThemeStore
	blobs    BlobStore
}

func (s *store) DB() *bolt.DB {
	return s.db
}

func (s *store) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *store) Sites() SiteStore       { return s.sites }
func (s *store) Pages() PageStore       { return s.pages }
func (s *store) Versions() VersionStore { return s.versions }
func (s *store) Themes() ThemeStore     { return s.themes }
func (s *store) Blobs() BlobStore       { return s.blobs }

func NewStore(rc *config.RuntimeConfig) (Store, error) {
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
		for b, n := range nestedBuckets {
			if _, err := tx.Bucket([]byte(b)).CreateBucketIfNotExists([]byte(n)); err != nil {
				return err
			}
		}
		return nil
	})

	blobs, err := NewBlobStore(rc, db)
	if err != nil {
		return nil, err
	}

	pageVersions := NewPageVersionIndex(db)

	return &store{
		db:       db,
		sites:    NewSiteStore(db),
		pages:    NewPageStore(db),
		versions: NewVersionStore(db, pageVersions),
		themes:   NewThemeStore(db),
		blobs:    blobs,
	}, nil
}
