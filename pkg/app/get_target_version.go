package app

import (
	"strings"

	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

func getTargetVersion(ctx *fiber.Ctx, svc Service, write bool) (*core.TargetVersion, error) {
	site, err := svc.GetSite(ctx.Context())
	if err != nil {
		return nil, err
	}
	versionHeader := ctx.GetReqHeaders()[HeaderPagebinVersion]

	if len(versionHeader) == 0 {
		if write {
			return core.NewNextTargetVersion(site.NextVersion), nil
		} else {
			return core.NewCurrentTargetVersion(site.Version), nil
		}
	}

	v := strings.TrimSpace(versionHeader[0])
	if len(v) == 0 {
		return nil, core.ErrInvalidVersion.New("%s header invalid. Specify either \"next\" or a version UID", HeaderPagebinVersion)
	}
	if v == "next" {
		return core.NewNextTargetVersion(site.NextVersion), nil
	}
	parsed, err := ulid.Parse(v)
	if err != nil {
		return nil, core.ErrInvalidVersion.New("%s header invalid. Specify either \"next\" or a version UID", HeaderPagebinVersion)
	}
	switch {
	case parsed == site.NextVersion:
		return core.NewNextTargetVersion(site.NextVersion), nil
	case parsed == site.Version:
		return core.NewCurrentTargetVersion(site.Version), nil
	default:
		if _, err := svc.GetVersion(ctx.Context(), parsed); err != nil {
			return nil, err
		}
		return core.NewTargetVersion(parsed), nil
	}
}
