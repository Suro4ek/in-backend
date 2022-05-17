package composites

import (
	"in-backend/internal/auth"
	"in-backend/internal/config"
	"in-backend/internal/handlers"
	"in-backend/internal/user"
	"in-backend/pkg/logging"
)

type AuthComposite struct {
	Handler handlers.HandlerAuth
}

func NewAuthComposite(logger *logging.Logger, Repository user.Repository, cfg *config.Config) (*AuthComposite, error) {
	handler := auth.NewHandler(logger, Repository, cfg)
	return &AuthComposite{
		handler,
	}, nil
}
