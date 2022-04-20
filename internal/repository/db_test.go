package repository

import (
	"context"
	"fmt"
	"github.com/IakimenkoD/xm-companies-service/internal/config"
	"github.com/IakimenkoD/xm-companies-service/internal/model"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/database"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/dataprovider"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/dataprovider/pg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"os"
	"testing"
)

// maybe there are more elegant solutions...
var companyStorage dataprovider.CompaniesStorage

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	log, _ := zap.NewProduction()

	testCfg := &config.Config{DB: config.DB{
		URL:          "postgres://root@localhost:5432/root?sslmode=disable",
		SchemaName:   "xm_test",
		MaxOpenConns: 2,
		MaxIdleConns: 2,
	}}

	testClient, err := database.NewClient(testCfg)
	if err != nil {
		return -1, fmt.Errorf("could not connect to database: %w", err)
	}
	companyStorage = pg.NewCompanyStorage(testClient, log)
	if err = testClient.Migrate(); err != nil {
		return -1, fmt.Errorf("could not migrate database: %w", err)
	}

	defer func() {
		_, _ = testClient.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", testCfg.DB.SchemaName))
		testClient.Close()
	}()

	return m.Run(), nil
}

func TestInsertCompany(t *testing.T) {
	b, err := companyStorage.Insert(context.TODO(), &model.Company{
		Name:    "Test1",
		Code:    "0101",
		Country: "UK",
		Website: "test01@tz.com",
		Phone:   "+6912345",
	})

	require.NoError(t, err)

	assert.Equal(t, int64(1), b)
}

func TestUpsertCompany(t *testing.T) {
	ctx := context.TODO()

	err := companyStorage.Update(ctx, &model.Company{
		ID:      1,
		Name:    "Test2",
		Code:    "10101",
		Country: "KU",
		Website: "test10@zt.com",
		Phone:   "+96543321",
	})

	require.NoError(t, err)

	f := dataprovider.NewCompanyFilter().ByIDs(1)
	c, err := companyStorage.GetByFilter(ctx, f)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, int64(1), c.ID)
	assert.Equal(t, "Test2", c.Name)
}

func TestDeleteCompany(t *testing.T) {
	err := companyStorage.DeleteByID(context.TODO(), 1)
	require.NoError(t, err)
	assert.NoError(t, err)

	f := dataprovider.NewCompanyFilter().ByIDs(1)
	c, err := companyStorage.GetByFilter(context.TODO(), f)
	assert.NoError(t, err)
	assert.Nil(t, c)
}
