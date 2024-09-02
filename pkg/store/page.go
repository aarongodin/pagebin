package store

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/oklog/ulid/v2"
	bolt "go.etcd.io/bbolt"
)

type PageStore interface {
	CreatePage(ctx context.Context, title string, path string, content ulid.ULID, templateName string, tags []string, excerpt string) (core.Page, error)
	GetPage(ctx context.Context, uid ulid.ULID) (core.Page, error)
}

type pageStore struct {
	db documentDB[core.Page]
}

func (s pageStore) CreatePage(ctx context.Context, title string, path string, content ulid.ULID, templateName string, tags []string, excerpt string) (core.Page, error) {
	page := core.Page{
		UID:          ulid.Make(),
		Title:        title,
		Path:         path,
		Content:      content,
		TemplateName: templateName,
		Tags:         tags,
		Excerpt:      excerpt,
	}
	if err := s.db.Save(ctx, bucketPages, page.UID.String(), page); err != nil {
		return core.Page{}, err
	}
	return page, nil
}

func (s pageStore) GetPage(ctx context.Context, uid ulid.ULID) (core.Page, error) {
	return s.db.One(ctx, bucketPages, uid.String())
}

func NewPageStore(db *bolt.DB) PageStore {
	return &pageStore{db: docDB[core.Page]{db}}
}
