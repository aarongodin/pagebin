package app

import (
	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/aymerick/raymond"
	"github.com/gofiber/fiber/v2"
)

type renderer struct {
	service Service
}

func (r renderer) render(ctx *fiber.Ctx) error {
	if isReservedPath(ctx.Path()) {
		return core.ErrReservedPath.NewWithNoMessage()
	}

	targetVersion, err := getTargetVersion(ctx, r.service, false)
	if err != nil {
		return err
	}

	pageUID, err := r.service.VersionManager().GetByPath(ctx.Context(), targetVersion, ctx.Path())
	if err != nil {
		return err
	}
	page, err := r.service.GetPage(ctx.Context(), pageUID)
	if err != nil {
		return err
	}
	content, err := r.service.ContentManager().Get(ctx.Context(), page.Content)
	if err != nil {
		return err
	}
	output, err := r.service.ThemeManager().Render(page.TemplateName, map[string]any{
		"content": raymond.SafeString(content),
	})
	if err != nil {
		return err
	}

	return ctx.SendString(string(output))
}

func NewRenderer(service Service) *renderer {
	return &renderer{service}
}
