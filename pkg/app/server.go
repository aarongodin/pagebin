package app

import (
	"github.com/aarongodin/pagebin/pkg/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	rc  *config.RuntimeConfig
	app *fiber.App
}

func (s *Server) Start() error {
	return s.app.Listen(s.rc.ServerAddr())
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
