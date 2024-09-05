package app

import (
	"context"
	"strings"

	"github.com/aarongodin/pagebin/pkg/config"
	"github.com/gofiber/fiber/v2"
)

const (
	HeaderPagebinVersion = "X-Pagebin-Version"
	pathAPI              = "/api"
)

var (
	reservedPaths = [...]string{pathAPI}
)

type Server struct {
	rc  *config.RuntimeConfig
	app *fiber.App
}

func (s *Server) Start() error {
	return s.app.Listen(s.rc.ServerAddr())
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

func NewServer(rc *config.RuntimeConfig, service Service) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	api := NewAdminAPI(service)
	api.Register(app)
	renderer := NewRenderer(service)
	app.Get("*", renderer.render)
	return &Server{rc, app}
}

func isReservedPath(path string) bool {
	for _, p := range reservedPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
