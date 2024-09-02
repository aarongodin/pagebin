package store

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/oklog/ulid/v2"
	bolt "go.etcd.io/bbolt"
)

type ThemeStore interface {
	CreateTheme(ctx context.Context, templates map[string]ulid.ULID, cssAssets []ulid.ULID, jsAssets []ulid.ULID) (core.Theme, error)
	GetTheme(ctx context.Context, uid ulid.ULID) (core.Theme, error)
}

type themeStore struct {
	db documentDB[core.Theme]
}

func (s themeStore) CreateTheme(ctx context.Context, templates map[string]ulid.ULID, cssAssets []ulid.ULID, jsAssets []ulid.ULID) (core.Theme, error) {
	theme := core.Theme{
		UID:       ulid.Make(),
		Templates: templates,
		CSSAssets: cssAssets,
		JSAssets:  jsAssets,
	}
	if err := s.db.Save(ctx, bucketThemes, theme.UID.String(), theme); err != nil {
		return core.Theme{}, err
	}
	return theme, nil
}

func (s themeStore) GetTheme(ctx context.Context, uid ulid.ULID) (core.Theme, error) {
	return s.db.One(ctx, bucketThemes, uid.String())
}

func NewThemeStore(db *bolt.DB) ThemeStore {
	return &themeStore{db: docDB[core.Theme]{db}}
}
