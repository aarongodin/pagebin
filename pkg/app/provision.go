package app

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/store"
	"github.com/joomcode/errorx"
	"github.com/oklog/ulid/v2"
)

var defaultSiteTitle = "New Site"
var defaultThemeTemplate = `
<!doctype html>
<html>
<head>
	<!-- pagebin:assets:css -->
</head>
<body>
	{{{ content }}}
	<!-- pagebin:assets:js -->
</body>
</html>
`
var defaultPage = `
<h1>Welcome to Pagebin</h1>
`

// Provision adds entities to the DB that are the minimum required for pagebin to function, such as a site, a theme, and a page.
func Provision(ctx context.Context, sites store.SiteStore, themes store.ThemeStore, pages store.PageStore, versions store.VersionStore, blobs store.BlobStore) error {
	shouldProvision := false
	_, err := sites.GetSite(ctx)
	if err != nil {
		if errorx.IsNotFound(err) {
			shouldProvision = true
		} else {
			return err
		}
	}
	if !shouldProvision {
		return nil
	}
	pageBlob, err := blobs.PutBytes(ctx, []byte(defaultPage))
	if err != nil {
		return err
	}
	page, err := pages.CreatePage(ctx, "Home", "/", pageBlob, "default", nil, "")
	if err != nil {
		return err
	}
	templateBlob, err := blobs.PutBytes(ctx, []byte(defaultThemeTemplate))
	if err != nil {
		return err
	}
	templates := map[string]ulid.ULID{
		"default": templateBlob,
	}
	theme, err := themes.CreateTheme(ctx, templates, nil, nil)
	if err != nil {
		return err
	}
	version, err := versions.CreateVersion(ctx, []ulid.ULID{page.UID}, theme.UID)
	if err != nil {
		return err
	}
	if _, err = sites.CreateSite(ctx, defaultSiteTitle, version.UID); err != nil {
		return err
	}
	return nil
}
