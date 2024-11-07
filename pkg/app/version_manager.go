package app

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/aarongodin/pagebin/pkg/store"
	"github.com/oklog/ulid/v2"
)

type VersionManager interface {
	GetByPath(ctx context.Context, targetVersion *core.TargetVersion, path string) (ulid.ULID, error)
	Load(ctx context.Context, currentUID ulid.ULID, nextUID ulid.ULID) error
	SetPage(previousPath string, path string, pageUID ulid.ULID) error
	UnsetPage(path string) error
}

type versionManager struct {
	current  *compiledVersion
	next     *compiledVersion
	versions store.VersionStore
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

func (m *versionManager) GetByPath(ctx context.Context, targetVersion *core.TargetVersion, path string) (ulid.ULID, error) {
	var targetCompiledVersion *compiledVersion
	switch {
	case targetVersion.IsCurrent():
		targetCompiledVersion = m.current
	case targetVersion.IsNext():
		targetCompiledVersion = m.next
	default:
		version, err := m.versions.GetVersion(ctx, targetVersion.UID())
		if err != nil {
			return ulid.ULID{}, err
		}
		pageUID, exists := version.Pages[path]
		if !exists {
			return ulid.ULID{}, core.ErrPageNotFound.NewWithNoMessage()
		}
		return pageUID, nil
	}

	if targetCompiledVersion == nil || targetCompiledVersion.index == nil {
		return ulid.ULID{}, core.ErrVersionNotCompiled.New("version not compiled")
	}
	uid, exists := targetCompiledVersion.Find(path)
	if !exists {
		return ulid.ULID{}, core.ErrPageNotFound.NewWithNoMessage()
	}
	return uid, nil
}

func (m *versionManager) Load(ctx context.Context, currentUID ulid.ULID, nextUID ulid.ULID) error {
	current, err := m.load(ctx, currentUID)
	if err != nil {
		return err
	}
	next, err := m.load(ctx, nextUID)
	if err != nil {
		return err
	}
	m.current = current
	m.next = next
	return nil
}

func (m *versionManager) load(ctx context.Context, uid ulid.ULID) (*compiledVersion, error) {
	version, err := m.versions.GetVersion(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &compiledVersion{
		uid:   uid,
		index: version.Pages,
	}, nil
}

func (m *versionManager) SetPage(previousPath string, path string, pageUID ulid.ULID) error {
	if m.next == nil || m.next.index == nil {
		return core.ErrVersionNotCompiled.New("next version not compiled")
	}
	if previousPath != "" {
		delete(m.next.index, previousPath)
	}
	m.next.index[path] = pageUID
	return nil
}

func (m *versionManager) UnsetPage(path string) error {
	if m.next == nil || m.next.index == nil {
		return core.ErrVersionNotCompiled.New("next version not compiled")
	}
	delete(m.next.index, path)
	return nil
}

func NewVersionManager(versions store.VersionStore) VersionManager {
	return &versionManager{
		versions: versions,
	}
}
