package app

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/core"
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
func Provision(ctx context.Context, store store.Store) error {
	shouldProvision := false
	_, err := store.Sites().GetSite(ctx)
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
	// TOOD: start tx here that goes across all provisioning
	pageBlob, err := store.Blobs().CreateBlob(ctx, []byte(defaultPage))
	if err != nil {
		return err
	}
	page, err := store.Pages().PutPage(ctx, nil, core.WritablePage{
		Title:        "Home",
		Path:         "/",
		TemplateName: "default",
	}, pageBlob.UID)
	if err != nil {
		return err
	}
	templateBlob, err := store.Blobs().CreateBlob(ctx, []byte(defaultThemeTemplate))
	if err != nil {
		return err
	}
	templates := map[string]ulid.ULID{
		"default": templateBlob.UID,
	}
	theme, err := store.Themes().CreateTheme(ctx, templates, nil, nil)
	if err != nil {
		return err
	}
	pages := map[string]ulid.ULID{
		page.Path: page.UID,
	}
	version, err := store.Versions().CreateVersion(ctx, pages, theme.UID)
	if err != nil {
		return err
	}
	nextVersion, err := store.Versions().Clone(ctx, version.UID)
	if err != nil {
		return err
	}
	if _, err = store.Sites().CreateSite(ctx, defaultSiteTitle, version.UID, nextVersion.UID); err != nil {
		return err
	}
	return nil
}
