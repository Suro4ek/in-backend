package composites

import (
	"in-backend/internal/handlers"
	"in-backend/internal/items"
	"in-backend/internal/items/db"
	"in-backend/pkg/logging"
	"in-backend/pkg/postgres"
)

type ItemComposite struct {
	Repository items.Repository
	Handler    handlers.Handler
}

func NewItemComposite(client *postgres.Client, logger *logging.Logger) (*ItemComposite, error) {
	repository := db.NewRepository(*client, logger)
	handler := items.NewHandler(repository, logger)
	return &ItemComposite{
		Repository: repository,
		Handler:    handler,
	}, nil
}
