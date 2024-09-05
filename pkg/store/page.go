package store

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/oklog/ulid/v2"
	bolt "go.etcd.io/bbolt"
)

type PageStore interface {
	PutPage(ctx context.Context, uid *ulid.ULID, write core.WritablePage, content ulid.ULID) (core.Page, error)
	GetPage(ctx context.Context, uid ulid.ULID) (core.Page, error)
}

type pageStore struct {
	db documentDB[core.Page]
}

func (s pageStore) PutPage(ctx context.Context, uid *ulid.ULID, write core.WritablePage, content ulid.ULID) (core.Page, error) {
	if uid == nil {
		newUID := ulid.Make()
		uid = &newUID
	}
	page := core.Page{
		UID:          *uid,
		Title:        write.Title,
		Path:         write.Path,
		Content:      content,
		TemplateName: write.TemplateName,
		Tags:         write.Tags,
		Excerpt:      write.Excerpt,
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
