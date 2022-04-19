package dataprovider

import (
	"context"
	"github.com/IakimenkoD/xm-companies-service/internal/model"
)

type CompaniesStore interface {
	GetByFilter(ctx context.Context, filter interface{}) (*model.Company, error)
	GetListByFilter(ctx context.Context, filter interface{}) ([]*model.Company, error)
	DeleteByFilter(ctx context.Context, filter interface{}) error

	Upsert(ctx context.Context, identityUserCode *model.Company) error
}
