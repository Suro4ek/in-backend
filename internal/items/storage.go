package items

import "context"

type Repository interface {
	GetOne(ctx context.Context, id string) (Item, error)
	GetAll(ctx context.Context) (i []Item, err error)
	Create(ctx context.Context, item *Item) error
	Update(ctx context.Context, item Item) error
	Delete(ctx context.Context, id string) error
}
