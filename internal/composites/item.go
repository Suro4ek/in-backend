package composites

import (
	"in-backend/internal/handlers"
	"in-backend/internal/items"
	"in-backend/internal/items/db"
	"in-backend/internal/user"
	"in-backend/pkg/logging"
	"in-backend/pkg/postgres"
)

type ItemComposite struct {
	Repository items.Repository
	Handler    handlers.HandlerAuth
}

func NewItemComposite(client *postgres.Client, userRepository user.Repository, logger *logging.Logger) (*ItemComposite, error) {
	repository := db.NewRepository(*client, logger)
	handler := items.NewHandler(repository, userRepository, logger)
	return &ItemComposite{
		Repository: repository,
		Handler:    handler,
	}, nil
}
