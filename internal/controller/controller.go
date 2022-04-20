package controller

import (
	"context"
	"github.com/IakimenkoD/xm-companies-service/internal/config"
	ierr "github.com/IakimenkoD/xm-companies-service/internal/errors"
	"github.com/IakimenkoD/xm-companies-service/internal/model"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/dataprovider"
)

// CompaniesService is a main controller for business logic.
type CompaniesService interface {
	CreateCompany(ctx context.Context, company *model.Company) (int64, error)
	GetCompanyByID(ctx context.Context, id int64) (*model.Company, error)
	GetCompanies(ctx context.Context, filter *dataprovider.CompanyFilter) ([]*model.Company, error)
	UpdateCompany(ctx context.Context, company *model.Company) error
	DeleteCompany(ctx context.Context, id int64) error

	//NotifyChanged
}

type Controller struct {
	config         *config.Config
	companyStorage dataprovider.CompaniesStorage
}

func (c Controller) GetCompanyByID(ctx context.Context, id int64) (*model.Company, error) {
	filter := dataprovider.NewCompanyFilter().ByIDs(id)
	company, err := c.companyStorage.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ierr.NotFound
	}
	return company, nil
}

func (c Controller) CreateCompany(ctx context.Context, company *model.Company) (id int64, err error) {
	return c.companyStorage.Insert(ctx, company)
}

func (c Controller) GetCompanies(ctx context.Context, filter *dataprovider.CompanyFilter) ([]*model.Company, error) {
	return c.companyStorage.GetListByFilter(ctx, filter)
}

func (c Controller) UpdateCompany(ctx context.Context, company *model.Company) error {
	return c.companyStorage.Update(ctx, company)
}

func (c Controller) DeleteCompany(ctx context.Context, id int64) error {
	filter := dataprovider.NewCompanyFilter().ByIDs(id)
	company, err := c.companyStorage.GetByFilter(ctx, filter)
	if err != nil {
		return err
	}
	if company == nil {
		return ierr.NotFound
	}
	return c.companyStorage.DeleteByID(ctx, id)
}

func NewCompaniesService(cfg *config.Config, companyStorage dataprovider.CompaniesStorage) CompaniesService {
	return &Controller{
		config:         cfg,
		companyStorage: companyStorage}
}
