package controller

import (
	"context"
	"github.com/IakimenkoD/xm-companies-service/internal/config"
	ierr "github.com/IakimenkoD/xm-companies-service/internal/errors"
	"github.com/IakimenkoD/xm-companies-service/internal/model"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/dataprovider"
	"github.com/IakimenkoD/xm-companies-service/internal/service"
)

//go:generate minimock -i CompaniesService -g -o controller_mock.go

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
	mq             service.MessageQueue
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
	if company == nil {
		return id, ierr.BadRequest
	}
	f := dataprovider.NewCompanyFilter().ByNames(company.Name).ByCodes(company.Code)
	duplicates, err := c.companyStorage.GetListByFilter(ctx, f)
	if err != nil {
		return id, err
	}
	if len(duplicates) > 0 {
		return id, ierr.CompanyExists
	}

	id, err = c.companyStorage.Insert(ctx, company)
	if err != nil {
		return id, err
	}
	company.ID = id
	if err = c.mq.NotifyCompanyUpdated(company); err != nil {
		return id, err
	}
	return id, nil
}

func (c Controller) GetCompanies(ctx context.Context, filter *dataprovider.CompanyFilter) ([]*model.Company, error) {
	return c.companyStorage.GetListByFilter(ctx, filter)
}

func (c Controller) UpdateCompany(ctx context.Context, company *model.Company) error {
	if company == nil {
		return ierr.BadRequest
	}

	f := dataprovider.NewCompanyFilter().ByIDs(company.ID)
	old, err := c.companyStorage.GetByFilter(ctx, f)
	if err != nil {
		return err
	}

	if old.Equal(company) {
		return nil
	}

	if err = c.companyStorage.Update(ctx, company); err != nil {
		return err
	}

	updated, err := c.companyStorage.GetByFilter(ctx, f)
	if err != nil {
		return err
	}

	if err = c.mq.NotifyCompanyUpdated(updated); err != nil {
		return err
	}
	return nil
}

func (c Controller) PatchCompany(ctx context.Context, company *model.Company) (*model.Company, error) {
	if company == nil {
		return nil, ierr.BadRequest
	}

	if err := c.companyStorage.Update(ctx, company); err != nil {
		return nil, err
	}

	f := dataprovider.NewCompanyFilter().ByIDs()
	updated, err := c.companyStorage.GetByFilter(ctx, f)
	if err != nil {
		return nil, err
	}
	if err = c.mq.NotifyCompanyUpdated(updated); err != nil {
		return nil, err
	}
	return updated, nil
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

func NewCompaniesService(cfg *config.Config,
	companyStorage dataprovider.CompaniesStorage,
	mq service.MessageQueue) CompaniesService {
	return &Controller{
		config:         cfg,
		companyStorage: companyStorage,
		mq:             mq,
	}
}
