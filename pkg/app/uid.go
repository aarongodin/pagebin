package app

import (
	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

func getUIDParam(ctx *fiber.Ctx, name string) (ulid.ULID, error) {
	raw := ctx.Params(name)
	if raw == "" {
		return ulid.ULID{}, core.ErrUIDRequired.NewWithNoMessage()
	}
	parsed, err := ulid.Parse(raw)
	if err != nil {
		return ulid.ULID{}, core.ErrUIDRequired.WrapWithNoMessage(err)
	}
	return parsed, nil
}
