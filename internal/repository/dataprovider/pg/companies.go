package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/IakimenkoD/xm-companies-service/internal/model"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/database"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/dataprovider"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"time"
)

func NewCompanyStorage(client *database.Client, logger *zap.Logger) dataprovider.CompaniesStorage {
	return &CompanyStore{
		db:     client,
		schema: client.SchemaName,
		log:    logger,
	}
}

type CompanyStore struct {
	db     *database.Client
	schema string
	log    *zap.Logger
}

func (s *CompanyStore) GetByFilter(ctx context.Context, filter *dataprovider.CompanyFilter) (*model.Company, error) {

	entities, err := s.GetListByFilter(ctx, filter)

	switch {
	case err != nil:
		return nil, err
	case len(entities) == 0:
		return nil, nil
	default:
		return entities[0], nil
	}
}

func (s *CompanyStore) GetListByFilter(ctx context.Context, filter *dataprovider.CompanyFilter) ([]*model.Company, error) {
	qb := sq.Select(
		"companies.id",
		"companies.name",
		"companies.code",
		"companies.country",
		"companies.website",
		"companies.phone",
		"companies.created_at",
		"companies.updated_at",
	).
		From(s.schema + ".companies").
		Where(getCompaniesCond(filter))

	query, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "creating sql query for getting companies by filter")
	}

	s.log.Debug("selecting company query SQL",
		zap.String("query", query),
		zap.Any("args", args))

	companies := []*model.Company{}
	if err = sqlx.SelectContext(ctx, s.db, &companies, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "selecting companies by filter from database with query %s", query)
	}

	return companies, nil
}

func (s *CompanyStore) Insert(ctx context.Context, company *model.Company) (id int64, err error) {
	query, args, err := sq.Insert(s.schema + ".companies").
		SetMap(map[string]interface{}{
			"name":       company.Name,
			"code":       company.Code,
			"country":    company.Country,
			"website":    company.Website,
			"phone":      company.Phone,
			"created_at": time.Now().UTC(),
		}).
		Suffix("RETURNING id;").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return id, errors.Wrap(err, "can't create query SQL for inserting company")
	}

	s.log.Debug("inserting company query SQL",
		zap.String("query", query),
		zap.Any("args", args))
	row := s.db.QueryRowxContext(ctx, query, args...)
	if err = row.Err(); err != nil {
		return id, errors.Wrap(err, "can't execute SQL query for inserting company")
	}

	err = row.Scan(&id)

	return id, errors.Wrap(err, "can't scan inserted company id")
}

func (s *CompanyStore) Update(ctx context.Context, company *model.Company) error {
	updates := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}

	if !emptyString(company.Name) {
		updates["name"] = company.Name
	}

	if !emptyString(company.Code) {
		updates["code"] = company.Code
	}

	if !emptyString(company.Country) {
		updates["country"] = company.Country
	}

	if !emptyString(company.Website) {
		updates["website"] = company.Website
	}

	if !emptyString(company.Phone) {
		updates["phone"] = company.Name
	}

	query, args, err := sq.Update(s.schema + ".companies").
		SetMap(updates).
		Where(sq.Eq{"id": company.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "creating sql query for updating company")
	}

	s.log.Debug("updating company query SQL",
		zap.String("query", query),
		zap.Any("args", args))

	_, err = s.db.ExecContext(ctx, query, args...)

	return errors.Wrap(err, "can't execute SQL query for updating company")
}

func (s *CompanyStore) DeleteByID(ctx context.Context, id int64) error {
	query, args, err := sq.Delete(s.schema + ".companies").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(err, "creating sql query for deleting company")
	}

	s.log.Debug("deleting company query SQL",
		zap.String("query", query),
		zap.Any("args", args))

	_, err = s.db.ExecContext(ctx, query, args...)

	return err
}

func getCompaniesCond(filter *dataprovider.CompanyFilter) sq.Sqlizer {
	eq := make(sq.Eq)
	neq := make(sq.NotEq)
	var cond sq.Sqlizer = sq.And{eq, neq}

	if len(filter.IDs) > 0 {
		eq["companies.id"] = filter.IDs
	}

	if len(filter.UserIDs) > 0 {
		eq["users.id"] = filter.UserIDs
	}

	if len(filter.Names) > 0 {
		eq["companies.name"] = filter.Names
	}

	if len(filter.Codes) > 0 {
		eq["companies.code"] = filter.Codes
	}

	if len(filter.Countries) > 0 {
		eq["companies.country"] = filter.Countries
	}

	if len(filter.WebSites) > 0 {
		eq["companies.web_site"] = filter.WebSites
	}

	if len(filter.Phones) > 0 {
		eq["companies.phone"] = filter.Phones
	}
	fmt.Println(cond)
	return cond
}

func emptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
