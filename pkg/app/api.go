package app

import (
	"github.com/gofiber/fiber/v2"
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

func (api adminAPI) Register(app *fiber.App) {
	grp := app.Group("/api/v1")
	grp.Get("/site", api.GetSite)
	grp.Patch("/site", api.UpdateSite)
}

func NewAdminAPI(service Service) *adminAPI {
	return &adminAPI{service}
}

type updateSiteBody struct {
	Title string `json:"title"`
}
