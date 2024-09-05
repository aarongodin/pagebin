package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	store, err := store.NewStore(rc)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init DB")
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init blob store")
	}
	svc, err := app.NewService(rc, store)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init service")
	}

	if err := app.Provision(ctx, store); err != nil {
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
	if err := svc.VersionManager().Load(ctx, site.Version, site.NextVersion); err != nil {
		log.Fatal().Err(err).Str("versionUID", site.Version.String()).Msg("failed compiling version on startup")
	}
	if err := svc.ThemeManager().Load(ctx, version.Theme); err != nil {
		log.Fatal().Err(err).Str("themeUID", version.Theme.String()).Msg("failed compiling theme on startup")
	}

	server := app.NewServer(rc, svc)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()
	log.Info().Int("port", rc.Port).Str("host", rc.Host).Msg("started http server")

	<-done
	log.Info().Msg("starting graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		if err := store.Close(ctx); err != nil {
			log.Err(err).Msg("failed to close DB")
		}
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Err(err).Msg("failed to shutdown http server")
	}

	log.Info().Msg("shutdown complete")

}
