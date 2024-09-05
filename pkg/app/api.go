package app

import (
	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type adminAPI struct {
	service Service
}

func (api adminAPI) GetSite(ctx *fiber.Ctx) error {
	site, err := api.service.GetSite(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(site)
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

func (api adminAPI) Register(app *fiber.App) {
	grp := app.Group("/api")
	grp.Get("/site", api.GetSite)
	grp.Patch("/site", api.UpdateSite)
	grp.Put("/pages", api.PutPage)
}

func NewAdminAPI(service Service) *adminAPI {
	return &adminAPI{service}
}

type updateSiteBody struct {
	Title string `json:"title"`
}

type pageBody struct {
	UID     *ulid.ULID
	Page    core.WritablePage `json:"page"`
	Content string            `json:"content"`
}
