package app

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/config"
	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/aarongodin/pagebin/pkg/store"
	"github.com/oklog/ulid/v2"
)

// Service contains the shared business logic for the app
type Service interface {
	GetSite(ctx context.Context) (core.Site, error)
	UpdateSite(ctx context.Context, title string) (core.Site, error)
	GetVersion(ctx context.Context, uid ulid.ULID) (core.Version, error)
	GetPage(ctx context.Context, uid ulid.ULID) (core.Page, error)
	VersionManager() VersionManager
	ThemeManager() ThemeManager
	ContentManager() ContentManager
}

type Svc struct {
	sites    store.SiteStore
	pages    store.PageStore
	versions store.VersionStore
	vm       VersionManager
	tm       ThemeManager
	cm       ContentManager
}

func (s *Svc) VersionManager() VersionManager {
	return s.vm
}

func (s *Svc) ThemeManager() ThemeManager {
	return s.tm
}

func (s *Svc) ContentManager() ContentManager {
	return s.cm
}

func (s *Svc) GetSite(ctx context.Context) (core.Site, error) {
	return s.sites.GetSite(ctx)
}

func (s *Svc) UpdateSite(ctx context.Context, title string) (core.Site, error) {
	return s.sites.UpdateSite(ctx, title)
}

func (s *Svc) GetPage(ctx context.Context, uid ulid.ULID) (core.Page, error) {
	return s.pages.GetPage(ctx, uid)
}

func (s *Svc) GetVersion(ctx context.Context, uid ulid.ULID) (core.Version, error) {
	return s.versions.GetVersion(ctx, uid)
}

func NewService(rc *config.RuntimeConfig, sites store.SiteStore, versions store.VersionStore, pages store.PageStore, themes store.ThemeStore, blobs store.BlobStore) (Service, error) {
	cm, err := NewCachedContentManager(rc, blobs)
	if err != nil {
		return nil, err
	}
	return &Svc{
		sites:    sites,
		pages:    pages,
		versions: versions,
		vm:       NewVersionManager(versions, pages),
		tm:       NewThemeManager(themes, blobs),
		cm:       cm,
	}, nil
}
