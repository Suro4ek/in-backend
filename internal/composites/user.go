package composites

import (
	"in-backend/internal/config"
	"in-backend/internal/handlers"
	"in-backend/internal/user"
	"in-backend/internal/user/db"
	"in-backend/pkg/logging"
	"in-backend/pkg/postgres"
)

type UserComposite struct {
	Repository user.Repository
	Handler    handlers.HandlerAuth
}

func NewUserComposite(client *postgres.Client, logger *logging.Logger, cfg *config.Config) (*UserComposite, error) {
	repository := db.NewRepository(*client, logger)
	handler := user.NewHandler(repository, logger, cfg)
	return &UserComposite{
		Repository: repository,
		Handler:    handler,
	}, nil
}
