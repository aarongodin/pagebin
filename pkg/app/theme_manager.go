package app

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/aarongodin/pagebin/pkg/store"
	"github.com/aymerick/raymond"
	"github.com/oklog/ulid/v2"
)

// ThemeManager controls how a theme is used during rendering
type ThemeManager interface {
	Render(templateName string, data any) ([]byte, error)
	Load(ctx context.Context, uid ulid.ULID) error
}

type themeManager struct {
	current *compiledTheme
	themes  store.ThemeStore
	blob    store.BlobStore
}

func (m themeManager) Render(templateName string, data any) ([]byte, error) {
	if m.current == nil {
		return nil, core.ErrThemeNotCompiled.NewWithNoMessage()
	}
	tpl, ok := m.current.templates[templateName]
	if !ok {
		return nil, core.ErrThemeTemplateNotFound.New("template \"%s\" not found for current theme", templateName)
	}
	out, err := tpl.Exec(data)
	if err != nil {
		return nil, core.ErrThemeTemplateExec.Wrap(err, "template \"%s\" failed to execute", templateName)
	}
	return []byte(out), nil
}

func (m *themeManager) Load(ctx context.Context, uid ulid.ULID) error {
	c := &compiledTheme{
		uid:       uid,
		templates: map[string]*raymond.Template{},
	}
	theme, err := m.themes.GetTheme(ctx, uid)
	if err != nil {
		return err
	}
	for templateName, templateUID := range theme.Templates {
		tplBlob, err := m.blob.GetBytes(ctx, templateUID)
		if err != nil {
			return err
		}
		tpl, err := raymond.Parse(string(tplBlob))
		if err != nil {
			return err
		}
		c.templates[templateName] = tpl
	}
	// TODO: CSS and JS asset caching
	m.current = c
	return nil
}

type compiledTheme struct {
	uid       ulid.ULID
	templates map[string]*raymond.Template
}

func NewThemeManager(themeStore store.ThemeStore, blobStore store.BlobStore) ThemeManager {
	return &themeManager{
		themes: themeStore,
		blob:   blobStore,
	}
}
