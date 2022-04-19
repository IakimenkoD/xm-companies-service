package dataprovider

import (
	"context"
	"github.com/IakimenkoD/xm-companies-service/internal/model"
)

type UsersStore interface {
	GetByFilter(ctx context.Context, filter interface{}) (*model.User, error)
	GetListByFilter(ctx context.Context, filter interface{}) ([]*model.User, error)
	DeleteByFilter(ctx context.Context, filter interface{}) error

	Upsert(ctx context.Context, identityUserCode *model.User) error
}
