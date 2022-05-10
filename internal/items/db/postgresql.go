package db

import (
	"context"
	"in-backend/internal/items"
	"in-backend/pkg/logging"
	"in-backend/pkg/postgres"
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

func NewRepository(client postgres.Client, logger *logging.Logger) items.Repository {
	client.DB.AutoMigrate(items.Item{})
	return &repository{
		client: client,
		logger: logger,
	}
}

func (r *repository) GetOne(ctx context.Context, id string) (items.Item, error) {
	r.logger.Tracef("Find item by id %s", id)
	var itm items.Item
	err := r.client.DB.Model(items.Item{}).First(&itm, id).Error
	if err != nil {
		return items.Item{}, nil
	}
	return itm, nil
}

func (r *repository) GetAll(ctx context.Context) (i []items.Item, err error) {
	r.logger.Trace("Find items")
	itemss := make([]items.Item, 1)
	err = r.client.DB.Find(&itemss).Error
	if err != nil {
		return nil, err
	}
	return itemss, nil
}

func (r *repository) Create(ctx context.Context, item *items.Item) error {
	r.logger.Tracef("Create item with serialNumber %s", item.SerialNumber)
	err := r.client.DB.Create(item).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, item items.Item) error {
	r.logger.Tracef("Update item id %d", item.ID)
	err := r.client.DB.Save(&item).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	r.logger.Tracef("Delete item with id %s", id)
	err := r.client.DB.Delete(&items.Item{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
