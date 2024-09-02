package main

import (
	"context"
	"os"

	"github.com/aarongodin/pagebin/pkg/app"
	"github.com/aarongodin/pagebin/pkg/config"
	"github.com/aarongodin/pagebin/pkg/store"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	rc, err := config.NewRuntimeConfig()
	if err != nil {
		log.Err(err).Msg("error parsing runtime config")
		return
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if rc.LogFormat == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	log.Info().Str("version", VERSION).Msg("pagebin")
	if rc.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Info().Str("lvl", zerolog.GlobalLevel().String()).Str("format", rc.LogFormat).Msg("logging config")

	db, err := store.NewDB(rc)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init DB")
	}
	siteStore := store.NewSiteStore(db)
	themeStore := store.NewThemeStore(db)
	pageStore := store.NewPageStore(db)
	versionStore := store.NewVersionStore(db)
	blobStore, err := store.NewBlobStore(rc)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init blob store")
	}
	svc, err := app.NewService(rc, siteStore, versionStore, pageStore, themeStore, blobStore)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init service")
	}

	if err := app.Provision(ctx, siteStore, themeStore, pageStore, versionStore, blobStore); err != nil {
		log.Fatal().Err(err).Msg("failed to provision app")
	}

	site, err := svc.GetSite(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed getting site on startup")
	}
	version, err := svc.GetVersion(ctx, site.Version)
	if err != nil {
		log.Fatal().Err(err).Str("versionUID", site.Version.String()).Msg("failed getting version on startup")
	}
	if err := svc.VersionManager().Load(ctx, site.Version); err != nil {
		log.Fatal().Err(err).Str("versionUID", site.Version.String()).Msg("failed compiling version on startup")
	}
	if err := svc.ThemeManager().Load(ctx, version.Theme); err != nil {
		log.Fatal().Err(err).Str("themeUID", version.Theme.String()).Msg("failed compiling theme on startup")
	}

	server := app.NewServer(rc, svc)
	server.Start()
}
