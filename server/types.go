package server

import (
	"github.com/calvinbui/teleport-sre-challenge/config"
	"github.com/calvinbui/teleport-sre-challenge/services"
)

type server struct {
	services *services.Services
	config   *config.Config
}

func New(svc *services.Services, conf *config.Config) *server {
	return &server{
		services: svc,
		config:   conf,
	}
}
