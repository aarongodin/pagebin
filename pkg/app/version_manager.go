package app

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/aarongodin/pagebin/pkg/store"
	"github.com/oklog/ulid/v2"
)

// VersionManager controls how a site version is used during rendering
type VersionManager interface {
	GetByPath(path string) (ulid.ULID, error)
	Load(ctx context.Context, uid ulid.ULID) error
}

type versionManager struct {
	current  *compiledVersion
	versions store.VersionStore
	pages    store.PageStore
}

// This could probably be placed in the core
type compiledVersion struct {
	uid   ulid.ULID
	index map[string]ulid.ULID // TODO: this is where I'd put an optimized index like a radix tree..... IF I HAD ONE
}

func (c compiledVersion) Find(value string) (ulid.ULID, bool) {
	uid, exists := c.index[value]
	return uid, exists
}

func (m *versionManager) GetByPath(path string) (ulid.ULID, error) {
	if m.current == nil || m.current.index == nil {
		return ulid.ULID{}, core.ErrVersionNotCompiled.NewWithNoMessage()
	}
	uid, exists := m.current.Find(path)
	if !exists {
		return ulid.ULID{}, core.ErrPageNotFound.NewWithNoMessage()
	}
	return uid, nil
}

func (m *versionManager) Load(ctx context.Context, uid ulid.ULID) error {
	c := &compiledVersion{
		uid:   uid,
		index: map[string]ulid.ULID{},
	}
	version, err := m.versions.GetVersion(ctx, uid)
	if err != nil {
		return err
	}
	// I'm stuck here, seems like I have to load every page, and maybe that's ok?
	// but perhaps I can just optimize it later. It doesn't seem like loading 10,000 json blobs would be that much of an issue though. Just will see a spike in compute
	for _, pageUID := range version.Pages {
		page, err := m.pages.GetPage(ctx, pageUID)
		if err != nil {
			return err
		}
		c.index[page.Path] = pageUID
	}
	m.current = c
	return nil
}

func NewVersionManager(versions store.VersionStore, pages store.PageStore) VersionManager {
	return &versionManager{
		versions: versions,
		pages:    pages,
	}
}
