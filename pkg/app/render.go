package app

import (
	"github.com/aymerick/raymond"
	"github.com/gofiber/fiber/v2"
)

type renderer struct {
	service Service
}

func (r renderer) render(ctx *fiber.Ctx) error {
	// get the page from the current version by request path
	pageUID, err := r.service.VersionManager().GetByPath(ctx.Path())
	if err != nil {
		return err
	}
	page, err := r.service.GetPage(ctx.Context(), pageUID)
	if err != nil {
		return err
	}

	// read the blob for the page if it exists
	content, err := r.service.ContentManager().Get(ctx.Context(), page.Content)
	if err != nil {
		return err
	}
	// pipe the blob through the theme / templating
	output, err := r.service.ThemeManager().Render(page.TemplateName, map[string]any{
		"content": raymond.SafeString(content),
	})
	if err != nil {
		return err
	}
	// send the response back to the context
	return ctx.SendString(string(output))
}

func NewRenderer(service Service) *renderer {
	return &renderer{service}
}
