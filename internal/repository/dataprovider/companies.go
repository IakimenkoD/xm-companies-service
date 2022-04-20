package dataprovider

import (
	"context"
	"github.com/IakimenkoD/xm-companies-service/internal/model"
)

type CompaniesStorage interface {
	GetByFilter(ctx context.Context, filter *CompanyFilter) (*model.Company, error)
	GetListByFilter(ctx context.Context, filter *CompanyFilter) ([]*model.Company, error)
	DeleteByID(ctx context.Context, id int64) error

	Insert(ctx context.Context, company *model.Company) (int64, error)
	Update(ctx context.Context, company *model.Company) error
}

// CompanyFilter is a filter for companies in storage.
type CompanyFilter struct {
	IDs       []int64
	UserIDs   []int64
	Names     []string
	Codes     []string
	Countries []string
	WebSites  []string
	Phones    []string
}

func NewCompanyFilter() *CompanyFilter {
	return &CompanyFilter{}
}

// ByIDs filters by xm.companies.id
func (f *CompanyFilter) ByIDs(ids ...int64) *CompanyFilter {
	f.IDs = ids
	return f
}

// ByUserIDs filters by xm.users.id
func (f *CompanyFilter) ByUserIDs(ids ...int64) *CompanyFilter {
	f.UserIDs = ids
	return f
}

// ByNames filters by xm.company.name
func (f *CompanyFilter) ByNames(names ...string) *CompanyFilter {
	f.Names = names
	return f
}

// ByNames filters by xm.company.name
func (f *CompanyFilter) ByCodes(codes ...string) *CompanyFilter {
	f.Codes = codes
	return f
}

// ByCountries filters by xm.company.country
func (f *CompanyFilter) ByCountries(countries ...string) *CompanyFilter {
	f.Countries = countries
	return f
}

// ByWebsites filters by xm.company.website
func (f *CompanyFilter) ByWebsites(websites ...string) *CompanyFilter {
	f.WebSites = websites
	return f
}

// ByPhones filters by xm.company.phone
func (f *CompanyFilter) ByPhones(phones ...string) *CompanyFilter {
	f.Phones = phones
	return f
}
