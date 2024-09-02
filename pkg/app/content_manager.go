package app

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/config"
	"github.com/aarongodin/pagebin/pkg/store"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/oklog/ulid/v2"
)

// ContentManager wraps logic for retrieving content for rendering
type ContentManager interface {
	Get(ctx context.Context, uid ulid.ULID) ([]byte, error)
}

type cachedContentManager struct {
	cache *lru.Cache[ulid.ULID, []byte]
	blob  store.BlobStore
}

func (m cachedContentManager) Get(ctx context.Context, uid ulid.ULID) ([]byte, error) {
	raw, cached := m.cache.Get(uid)
	if cached {
		return raw, nil
	}
	raw, err := m.blob.GetBytes(ctx, uid)
	if err != nil {
		return nil, err
	}
	m.cache.Add(uid, raw)
	return raw, nil
}

// NewCachedContentManager creates a content manager backed by an LRU cache. Configure the cache size and options through runtime config.
func NewCachedContentManager(rc *config.RuntimeConfig, blob store.BlobStore) (ContentManager, error) {
	cache, err := lru.New[ulid.ULID, []byte](rc.ContentCacheSize)
	if err != nil {
		return nil, err
	}
	return &cachedContentManager{cache, blob}, nil
}
