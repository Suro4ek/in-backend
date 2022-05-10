package user

import "context"

type Repository interface {
	GetOne(ctx context.Context, id string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetAll(ctx context.Context) (u []User, err error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
}
