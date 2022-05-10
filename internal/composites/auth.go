package composites

import (
	"in-backend/internal/auth"
	"in-backend/internal/handlers"
	"in-backend/internal/user"
	"in-backend/pkg/logging"
)

type AuthComposite struct {
	Handler handlers.Handler
}

func NewAuthComposite(logger *logging.Logger, Repository user.Repository) (*AuthComposite, error) {
	handler := auth.NewHandler(logger, Repository)
	return &AuthComposite{
		handler,
	}, nil
}
