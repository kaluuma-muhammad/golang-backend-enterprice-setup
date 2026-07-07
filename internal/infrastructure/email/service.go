package email

import (
	"github.com/go-api/internal/shared/config"
)

type Service struct {
	cfg      *config.Config
	provider Provider
	renderer *Renderer
}

func New(cfg *config.Config, provider Provider, renderer *Renderer) *Service {
	return &Service{
		cfg:      cfg,
		provider: provider,
		renderer: renderer,
	}
}
