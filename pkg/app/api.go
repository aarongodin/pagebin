package app

import (
	"net/http"

	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

func NewAdminAPI(service Service) *adminAPI {
	return &adminAPI{service}
}

type adminAPI struct {
	service Service
}

func (api adminAPI) Register(app *fiber.App) {
	grp := app.Group("/api")

	grp.Get("/site", api.GetSite)
	grp.Patch("/site", api.UpdateSite)

	// grp.Get("/pages", api.GetPages)
	grp.Get("/page/:uid", api.GetPage)
	grp.Put("/pages", api.PutPage)
	grp.Delete("/pages/:uid", api.DeletePage)

	grp.Get("/versions", api.GetVersions)
	grp.Get("/versions/:uid", api.GetVersion)

	grp.Get("/theme/:uid", api.GetTheme)
	grp.Get("/theme/:uid/template/:name", api.GetThemeTemplate)
	grp.Get("/theme/:uid/asset/:uid", api.GetThemeAsset)
}

func (api adminAPI) GetSite(ctx *fiber.Ctx) error {
	site, err := api.service.GetSite(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(site)
}

// func (api adminAPI) GetPages(ctx *fiber.Ctx) error {
// 	pages, cursor, err := api.service.GetPages(ctx.Context())
// 	if err != nil {
// 		return err
// 	}
// 	return ctx.JSON(paginated[[]core.Page]{
// 		Cursor: cursor,
// 		Items:  pages,
// 	})
// }

func (api adminAPI) GetPage(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotImplemented)
}

func (api adminAPI) GetVersions(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotImplemented)
}

func (api adminAPI) GetVersion(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotImplemented)
}

func (api adminAPI) GetTheme(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotImplemented)
}

func (api adminAPI) GetThemeAsset(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotImplemented)
}

func (api adminAPI) GetThemeTemplate(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotImplemented)
}

func (api adminAPI) UpdateSite(ctx *fiber.Ctx) error {
	b := updateSiteBody{}
	if err := ctx.BodyParser(&b); err != nil {
		return err
	}
	site, err := api.service.UpdateSite(ctx.Context(), b.Title)
	if err != nil {
		return err
	}
	return ctx.JSON(site)
}

func (api adminAPI) PutPage(ctx *fiber.Ctx) error {
	b := pageBody{}
	if err := ctx.BodyParser(&b); err != nil {
		return err
	}
	page, err := api.service.PutPage(ctx.Context(), b.UID, b.Page, []byte(b.Content))
	if err != nil {
		return err
	}
	return ctx.JSON(page)
}

func (api adminAPI) DeletePage(ctx *fiber.Ctx) error {
	uid, err := getUIDParam(ctx, "uid")
	if err != nil {
		return err
	}
	if err := api.service.DeletePage(ctx.Context(), uid); err != nil {
		return err
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func (api adminAPI) PutThemeTemplate(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotImplemented)
}

func (api adminAPI) DeleteThemeTemplate(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotImplemented)
}

func (api adminAPI) PutThemeAsset(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotImplemented)
}

func (api adminAPI) DeleteThemeAsset(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusNotImplemented)
}

type updateSiteBody struct {
	Title string `json:"title"`
}

type pageBody struct {
	UID     *ulid.ULID
	Page    core.WritablePage `json:"page"`
	Content string            `json:"content"`
}

type paginated[T []core.Page] struct {
	Cursor *ulid.ULID `json:"cursor"`
	Items  T          `json:"items"`
}
