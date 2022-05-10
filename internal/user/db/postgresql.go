package db

import (
	"context"
	"in-backend/internal/user"
	"in-backend/pkg/logging"
	"in-backend/pkg/postgres"
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

func NewRepository(client postgres.Client, logger *logging.Logger) user.Repository {
	client.DB.AutoMigrate(user.User{})
	return &repository{
		client: client,
		logger: logger,
	}
}

func (r *repository) GetOne(ctx context.Context, id string) (user.User, error) {
	r.logger.Tracef("Find user by id %s", id)
	var usr user.User
	err := r.client.DB.Model(user.User{}).First(&usr, id).Error
	if err != nil {
		return user.User{}, err
	}
	return usr, err
}

func (r *repository) GetByUsername(ctx context.Context, username string) (user.User, error) {
	r.logger.Tracef("Find user by username %s", username)
	var usr user.User
	err := r.client.DB.Model(user.User{}).First(&usr, "username = ?", username).Error
	if err != nil {
		return user.User{}, err
	}
	return usr, err
}

func (r *repository) GetAll(ctx context.Context) (u []user.User, err error) {
	r.logger.Trace("Find users")
	users := make([]user.User, 1)
	err = r.client.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) Create(ctx context.Context, user *user.User) error {
	r.logger.Tracef("Create user with username %s", user.Username)
	err := r.client.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, user user.User) error {
	r.logger.Tracef("Update user id %d", user.ID)
	err := r.client.DB.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	r.logger.Tracef("Delete user with id %s", id)
	err := r.client.DB.Delete(&user.User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
