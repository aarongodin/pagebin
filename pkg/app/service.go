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
	PutPage(ctx context.Context, uid *ulid.ULID, page core.WritablePage, content []byte) (created core.Page, txErr error)
	VersionManager() VersionManager
	ThemeManager() ThemeManager
	ContentManager() ContentManager
}

type Svc struct {
	store store.Store
	vm    VersionManager
	tm    ThemeManager
	cm    ContentManager
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
	return s.store.Sites().GetSite(ctx)
}

func (s *Svc) UpdateSite(ctx context.Context, title string) (core.Site, error) {
	return s.store.Sites().UpdateSite(ctx, title)
}

func (s *Svc) GetPage(ctx context.Context, uid ulid.ULID) (core.Page, error) {
	return s.store.Pages().GetPage(ctx, uid)
}

func (s *Svc) GetVersion(ctx context.Context, uid ulid.ULID) (core.Version, error) {
	return s.store.Versions().GetVersion(ctx, uid)
}

func (s *Svc) PutPage(ctx context.Context, uid *ulid.ULID, write core.WritablePage, content []byte) (created core.Page, txErr error) {
	ctx, err := s.store.StartTx(ctx, true)
	if err != nil {
		return created, err
	}
	defer func() {
		txErr = s.store.EndTx(ctx, txErr)
	}()

	var page *core.Page
	if uid != nil {
		p, err := s.GetPage(ctx, *uid)
		if err != nil {
			return created, err
		}
		page = &p
	}

	previousPath := ""
	if page != nil {
		previousPath = page.Path
	}

	if page == nil {
		p, err := s.createPage(ctx, write, content)
		if err != nil {
			return created, err
		}
		page = &p
	} else {
		p, err := s.updatePage(ctx, *page, write, content)
		if err != nil {
			return created, err
		}
		page = &p
	}

	if err := s.VersionManager().SetPage(previousPath, write.Path, page.UID); err != nil {
		return created, err
	}

	return *page, nil
}

func (s *Svc) createPage(ctx context.Context, write core.WritablePage, content []byte) (core.Page, error) {
	contentBlob, err := s.store.Blobs().CreateBlob(ctx, content)
	if err != nil {
		return core.Page{}, err
	}
	site, err := s.GetSite(ctx)
	if err != nil {
		return core.Page{}, err
	}
	page, err := s.store.Pages().PutPage(ctx, nil, write, contentBlob.UID)
	if err != nil {
		return core.Page{}, err
	}
	if _, err = s.store.Versions().SetPage(ctx, site.NextVersion, "", page.Path, page.UID); err != nil {
		return page, err
	}
	return page, nil
}

func (s *Svc) updatePage(ctx context.Context, current core.Page, write core.WritablePage, content []byte) (core.Page, error) {
	createPage := false
	site, err := s.GetSite(ctx)
	if err != nil {
		return core.Page{}, err
	}
	versions, err := s.store.Versions().GetPageVersions(ctx, current.UID)
	if err != nil {
		return core.Page{}, err
	}
	switch versions.Cardinality() {
	case 0:
		return core.Page{}, core.ErrUnknown.New("expected page %s to belong to at least one version", current.UID.String())
	case 1:
		v := versions.ToSlice()[0]
		if v != site.NextVersion {
			createPage = true
		}
	default:
		createPage = true
	}

	if createPage {
		if versions.Contains(site.NextVersion) {
			// this is where I would reconcile that the next version should not have the old UID
			// probably do soemthing like:
			// s.store.Versions().RemovePage(ctx, site.NextVersion, current.UID)
		}
		return s.createPage(ctx, write, content)
	}

	blob, err := s.store.Blobs().GetBlob(ctx, current.Content)
	if err != nil {
		return core.Page{}, err
	}
	if eq, err := blob.Equal(content); err != nil {
		return core.Page{}, err
	} else if !eq {
		if blob, err = s.store.Blobs().UpdateBlob(ctx, blob.UID, content); err != nil {
			return core.Page{}, err
		}
	}

	page, err := s.store.Pages().PutPage(ctx, &current.UID, write, blob.UID)
	if err != nil {
		return core.Page{}, err
	}

	return page, nil
}

// func (s *Svc) CreatePage(ctx context.Context, title string, path string, templateName string, tags []string, excerpt string, content []byte) (page core.Page, txErr error) {
// 	// TODO: verify template exists

// 	ctx, err := s.store.StartTx(ctx, true)
// 	if err != nil {
// 		return page, err
// 	}
// 	defer func() {
// 		txErr = s.store.EndTx(ctx, txErr)
// 	}()

// 	contentBlob, err := s.store.Blobs().CreateBlob(ctx, content)
// 	if err != nil {
// 		return page, err
// 	}

// 	site, txErr := s.GetSite(ctx)
// 	if txErr != nil {
// 		return page, txErr
// 	}
// 	page, txErr = s.store.Pages().PutPage(ctx, nil, title, path, contentBlob.UID, templateName, tags, excerpt)
// 	if txErr != nil {
// 		return page, txErr
// 	}
// 	if _, txErr = s.store.Versions().SetPage(ctx, site.NextVersion, "", page.Path, page.UID); txErr != nil {
// 		return page, txErr
// 	}
// 	if txErr := s.VersionManager().SetPage("", path, page.UID); txErr != nil {
// 		return page, txErr
// 	}
// 	return page, nil
// }

// func (s *Svc) UpdatePage(ctx context.Context, uid ulid.ULID, title string, path string, templateName string, tags []string, excerpt string, content []byte) (page core.Page, contentUpdated bool, txErr error) {
// 	// TODO: verify template exists

// 	ctx, err := s.store.StartTx(ctx, true)
// 	if err != nil {
// 		return page, false, err
// 	}
// 	defer func() {
// 		txErr = s.store.EndTx(ctx, txErr)
// 	}()

// 	site, txErr := s.GetSite(ctx)
// 	if txErr != nil {
// 		return page, false, txErr
// 	}

// 	page, txErr = s.store.Pages().GetPage(ctx, uid)
// 	if txErr != nil {
// 		return page, false, txErr
// 	}
// 	previousPath := page.Path

// 	if path != previousPath {
// 		if _, txErr = s.store.Versions().SetPage(ctx, site.NextVersion, previousPath, path, page.UID); txErr != nil {
// 			return page, false, txErr
// 		}
// 		if txErr := s.VersionManager().SetPage(previousPath, path, page.UID); txErr != nil {
// 			return page, false, txErr
// 		}
// 	}

// 	contentBlob, txErr := s.store.Blobs().GetBlob(ctx, page.Content)
// 	if txErr != nil {
// 		return page, false, txErr
// 	}
// 	contentBlobEq, txErr := contentBlob.Equal(content)
// 	if txErr != nil {
// 		return page, false, txErr
// 	}
// 	if !contentBlobEq {
// 		if _, txErr := s.store.Blobs().UpdateBlob(ctx, page.Content, content); txErr != nil {
// 			return page, false, txErr
// 		}
// 		contentUpdated = true
// 	}

// 	page, txErr = s.store.Pages().PutPage(ctx, &uid, title, path, page.Content, templateName, tags, excerpt)
// 	if err != nil {
// 		return page, false, txErr
// 	}

// 	return page, contentUpdated, nil
// }

func NewService(rc *config.RuntimeConfig, s store.Store) (Service, error) {
	cm, err := NewCachedContentManager(rc, s.Blobs())
	if err != nil {
		return nil, err
	}
	return &Svc{
		store: s,
		vm:    NewVersionManager(s.Versions()),
		tm:    NewThemeManager(s.Themes(), s.Blobs()),
		cm:    cm,
	}, nil
}
