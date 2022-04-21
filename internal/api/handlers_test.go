package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IakimenkoD/xm-companies-service/internal/config"
	"github.com/IakimenkoD/xm-companies-service/internal/controller"
	ierr "github.com/IakimenkoD/xm-companies-service/internal/errors"
	"github.com/IakimenkoD/xm-companies-service/internal/model"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/database"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/dataprovider"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/dataprovider/pg"
	"github.com/IakimenkoD/xm-companies-service/internal/service"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	cyLocation = "1.1.1.1"
	usLocation = "8.8.8.8"
	bsLocation = "10.10.10.10."

	testingToken = "dGVzdCBjYXNlIHJlcXVpcmVkIHRva2Vu"

	companiesURL = "/companies"
)

func TestCreateCompanies(t *testing.T) {

	tt := []testCase{
		{
			name:           "fail: auth required",
			path:           companiesURL,
			method:         http.MethodPost,
			prepareRequest: prepareRequest(`{"name": "my company","code": "1235","country": "CY","website": "example.com","phone": "+79991123123"}`, cyLocation),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "fail: all fields required",
			path:           companiesURL,
			method:         http.MethodPost,
			token:          testingToken,
			prepareRequest: prepareRequest(`{"name": "my company","code": "1235","website": "example.com","phone": "+79991123123"}`, cyLocation),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "country: Invalid param\n",
		},
		{
			name:           "fail: wrong location",
			path:           companiesURL,
			method:         http.MethodPost,
			token:          testingToken,
			prepareRequest: prepareRequest(`{"name": "my company","code": "1235","website": "example.com","phone": "+79991123123"}`, usLocation),
			expectedStatus: http.StatusForbidden,
			expectedBody:   "your location is not allowed\n",
		},
		{
			name:           "fail: invalid location",
			path:           companiesURL,
			method:         http.MethodPost,
			token:          testingToken,
			prepareRequest: prepareRequest(`{"name": "my company","code": "1235","website": "example.com","phone": "+79991123123"}`, bsLocation),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "unexpected response from external service\n",
		},
		{
			name:           "success",
			path:           companiesURL,
			method:         http.MethodPost,
			token:          testingToken,
			prepareRequest: prepareRequest(`{"name": "my company","code": "1235" ,"country":"CY","website": "example.com","phone": "+79991123123"}`, cyLocation),
			expectedStatus: http.StatusCreated,
		}}
	checkTestCases(t, tt)
}

func TestGetCompanies(t *testing.T) {
	//should be called in first test case
	prepareDB := func(_ *testing.T, db *store) {
		db.client.MustExec(`INSERT INTO ` + db.client.SchemaName + `.companies` +
			`  ( id,        name,        code, country,        website,     phone) VALUES` +
			`  ( 11,   'testOne',      '1111',    'cy',   'testone.cy',   '+001234')` +
			`, ( 12,   'testTwo',      '2222',    'uk',   'testtwo.uk',   '+002345')` +
			`, ( 13, 'testThree',      '3333',    'bg', 'testthree.bg',   '+003456')` +
			`, ( 14,  'testFour',      '4444',    'cy',  'testfour.cy',   '+004567')` +
			`;`)
	}

	tt := []testCase{
		{
			name:           "get by id",
			path:           companiesURL + "/14",
			method:         http.MethodGet,
			prepareDB:      prepareDB,
			prepareRequest: prepareRequest(nil, cyLocation),
			expectedStatus: http.StatusOK,
			afterTest: func(t *testing.T, resp *http.Response) {
				var companies []*model.Company
				if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
					t.Fatalf("could not decode response body: %+v", err)
				}
				if assert.Len(t, companies, 1) {
					assert.Equal(t, "testFour", companies[0].Name)
				}
			},
		},
		{
			name:           "get by single id",
			path:           companiesURL + "?ids=11",
			method:         http.MethodGet,
			prepareRequest: prepareRequest(nil, cyLocation),
			expectedStatus: http.StatusOK,
			afterTest: func(t *testing.T, resp *http.Response) {
				var companies []*model.Company
				if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
					t.Fatalf("could not decode response body: %+v", err)
				}
				if assert.Len(t, companies, 1) {
					assert.Equal(t, "testOne", companies[0].Name)
				}
			},
		},
		{
			name:           "get by couple ids",
			path:           companiesURL + "?ids=11,12",
			method:         http.MethodGet,
			prepareRequest: prepareRequest(nil, cyLocation),
			expectedStatus: http.StatusOK,
			afterTest: func(t *testing.T, resp *http.Response) {
				var companies []*model.Company
				if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
					t.Fatalf("could not decode response body: %+v", err)
				}
				if assert.Len(t, companies, 2) {
					assert.Equal(t, "testOne", companies[0].Name)
					assert.Equal(t, "testTwo", companies[1].Name)
				}
			},
		},
		{
			name:           "get by name",
			path:           companiesURL + "?names=testThree",
			method:         http.MethodGet,
			prepareRequest: prepareRequest(nil, cyLocation),
			expectedStatus: http.StatusOK,
			afterTest: func(t *testing.T, resp *http.Response) {
				var companies []*model.Company
				if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
					t.Fatalf("could not decode response body: %+v", err)
				}
				if assert.Len(t, companies, 1) {
					assert.EqualValues(t, 13, companies[0].ID)
				}
			},
		},
		{
			name:           "get by code",
			path:           companiesURL + "?codes=4444",
			method:         http.MethodGet,
			prepareRequest: prepareRequest(nil, cyLocation),
			expectedStatus: http.StatusOK,
			afterTest: func(t *testing.T, resp *http.Response) {
				var companies []*model.Company
				if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
					t.Fatalf("could not decode response body: %+v", err)
				}
				if assert.Len(t, companies, 1) {
					assert.EqualValues(t, 14, companies[0].ID)
				}
			},
		},
		{
			name:           "get by countries, case insensitive",
			path:           companiesURL + "?countries=cY",
			method:         http.MethodGet,
			prepareRequest: prepareRequest(nil, cyLocation),
			expectedStatus: http.StatusOK,
			afterTest: func(t *testing.T, resp *http.Response) {
				var companies []*model.Company
				if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
					t.Fatalf("could not decode response body: %+v", err)
				}
				if assert.Len(t, companies, 2) {
					assert.EqualValues(t, 11, companies[0].ID)
					assert.EqualValues(t, 14, companies[1].ID)
				}
			},
		},
		{
			name:           "get by website",
			path:           companiesURL + "?websites=testthree.bg",
			method:         http.MethodGet,
			prepareRequest: prepareRequest(nil, cyLocation),
			expectedStatus: http.StatusOK,
			afterTest: func(t *testing.T, resp *http.Response) {
				var companies []*model.Company
				if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
					t.Fatalf("could not decode response body: %+v", err)
				}
				if assert.Len(t, companies, 1) {
					assert.EqualValues(t, 13, companies[0].ID)
				}
			},
		},
		{
			name:           "get by phone",
			path:           companiesURL + "?phones=+004567",
			method:         http.MethodGet,
			prepareRequest: prepareRequest(nil, cyLocation),
			expectedStatus: http.StatusOK,
			afterTest: func(t *testing.T, resp *http.Response) {
				var companies []*model.Company
				if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
					t.Fatalf("could not decode response body: %+v", err)
				}
				if assert.Len(t, companies, 1) {
					assert.EqualValues(t, 14, companies[0].ID)
				}
			},
		},
	}
	checkTestCases(t, tt)
}

func TestUpdateCompanies(t *testing.T) {
	//should be called in first test case
	prepareDB := func(_ *testing.T, db *store) {
		db.client.MustExec(`INSERT INTO ` + db.client.SchemaName + `.companies` +
			`  ( id,        name,        code, country,        website,     phone) VALUES` +
			`  ( 11,   'testOne',      '1111',    'cy',   'testone.cy',   '+001234')` +
			`, ( 12,   'testTwo',      '2222',    'uk',   'testtwo.uk',   '+002345')` +
			`, ( 13, 'testThree',      '3333',    'bg', 'testthree.bg',   '+003456')` +
			`, ( 14,  'testFour',      '4444',    'cy',  'testfour.cy',   '+004567')` +
			`;`)
	}
	tt := []testCase{
		{
			name:           "success: no auth and location check required",
			path:           companiesURL + "/11",
			method:         http.MethodPut,
			prepareDB:      prepareDB,
			prepareRequest: prepareRequest(`{"name": "my company","code": "1235","country": "CY","website": "example.com","phone": "+79991123123"}`, ""),
			expectedStatus: http.StatusNoContent,
			checkDB: func(t *testing.T, stores *store) {
				f := dataprovider.NewCompanyFilter().ByIDs(11)
				company, err := stores.companyStorage.GetByFilter(context.Background(), f)
				assert.NoError(t, err)
				assert.NotNil(t, company)
				assert.EqualValues(t, "1235", company.Code)
			},
		},
		{
			name:           "fail: all fields required",
			path:           companiesURL + "/12",
			method:         http.MethodPut,
			token:          testingToken,
			prepareRequest: prepareRequest(`{"name": "my company", "website": "example.com","phone": "+79991123123"}`, cyLocation),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "code: Invalid param\n",
		},
		{
			name:           "success: partial update",
			path:           companiesURL + "/13",
			method:         http.MethodPatch,
			token:          testingToken,
			prepareRequest: prepareRequest(`{"name": "Meta","website": "google.com"}`, usLocation),
			expectedStatus: http.StatusOK,
			afterTest: func(t *testing.T, resp *http.Response) {
				company := model.Company{}
				if err := json.NewDecoder(resp.Body).Decode(&company); err != nil {
					t.Fatalf("could not decode response body: %+v", err)
				}
				assert.EqualValues(t, 13, company.ID)
				assert.EqualValues(t, "Meta", company.Name)
				assert.EqualValues(t, "google.com", company.Website)
			},
		},
	}
	checkTestCases(t, tt)
}

func TestDeleteCompanies(t *testing.T) {
	tt := []testCase{
		{
			name:           "fail: auth required",
			path:           companiesURL + "/11",
			method:         http.MethodDelete,
			prepareRequest: prepareRequest(nil, cyLocation),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "fail: wrong location",
			path:           companiesURL + "/11",
			method:         http.MethodDelete,
			token:          testingToken,
			prepareRequest: prepareRequest(nil, usLocation),
			expectedStatus: http.StatusForbidden,
			expectedBody:   "your location is not allowed\n",
		},
		{
			name:           "fail: company not found",
			path:           companiesURL + "/9",
			token:          testingToken,
			method:         http.MethodDelete,
			prepareRequest: prepareRequest(nil, cyLocation),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Company not found\n",
		},
		{
			name:   "success",
			path:   companiesURL + "/11",
			token:  testingToken,
			method: http.MethodDelete,
			prepareDB: func(_ *testing.T, db *store) {
				db.client.MustExec(`INSERT INTO ` + db.client.SchemaName + `.companies` +
					`  ( id,        name,        code, country,        website,     phone) VALUES` +
					`  ( 11,   'testOne',      '1111',    'CY',   'testone.cy',   '+001234')` +
					`;`)
			},
			prepareRequest: prepareRequest(nil, cyLocation),
			expectedStatus: http.StatusNoContent,
			checkDB: func(t *testing.T, stores *store) {
				f := dataprovider.NewCompanyFilter().ByIDs(11)
				company, err := stores.companyStorage.GetByFilter(context.Background(), f)
				assert.NoError(t, err)
				assert.Nil(t, company)
			},
		},
	}
	checkTestCases(t, tt)
}

func checkTestCases(t *testing.T, tt []testCase) {
	logger, _ := zap.NewDevelopment()
	defaultConf, _ := config.New("", logger)
	defaultConf.DB.SchemaName = "xm_test"

	dbClient, err := database.NewClient(defaultConf)
	if err != nil {
		panic(err.Error())
	}
	if err = dbClient.Migrate(); err != nil {
		panic(err)
	}
	storage := pg.NewCompanyStorage(dbClient, logger)

	mqMock := configureMqMock(service.NewMessageQueueMock(t))
	companiesService := controller.NewCompaniesService(defaultConf, storage, mqMock)

	srv, err := NewServer(defaultConf, companiesService, configureIpCheckerMock(service.NewIpCheckerMock(t)))
	if err != nil {
		panic(err)
	}

	store := &store{
		client:         dbClient,
		companyStorage: storage,
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			h := httptest.NewServer(srv.Handler)
			defer h.Close()

			if tc.prepareDB != nil {
				tc.prepareDB(t, store)
			}

			// prepare request
			baseURL := "/api/v1"

			URL := h.URL + baseURL + tc.path

			method := http.MethodGet
			if tc.method != "" {
				method = tc.method
			}
			fmt.Printf("Checking %s %q (%s)\n", method, URL, tc.name)
			// to add data in request body use prepareRequest. ioutil.NopCloser(io.Reader)
			req, err := http.NewRequestWithContext(context.Background(), method, URL, nil)
			assert.NoErrorf(t, err, "could not create request for %q, %v", URL, err)

			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			if tc.prepareRequest != nil {
				tc.prepareRequest(req)
			}

			// do request
			resp, err := http.DefaultClient.Do(req)
			// if err != nil and resp = nil => maybe some mocks are not configured
			assert.NoErrorf(t, err, "can't do request")
			assert.NotNil(t, resp)
			defer func(rc io.ReadCloser) {
				err = rc.Close()
				assert.NoError(t, err)
			}(resp.Body)

			// check response status
			status := http.StatusOK
			if tc.expectedStatus != 0 {
				status = tc.expectedStatus
			}

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("can't read response body: %+v", err)
			}

			// reset the response body to the original unread state
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			assert.Equal(t, status, resp.StatusCode)

			if tc.afterTest != nil {
				tc.afterTest(t, resp)
			}

			if tc.checkDB != nil {
				tc.checkDB(t, store)
			}

			// check response body
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoErrorf(t, err, "could not read response body, %v", err)
			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, string(body))
			}
			if t.Failed() {
				var prettyJSON bytes.Buffer
				if err = json.Indent(&prettyJSON, bodyBytes, "", "  "); err != nil {
					t.Logf("can't format JSON: %+v", err)
					t.Logf("raw response body:\n%s", string(body))
				} else {
					t.Logf("raw pretty JSON response:\n%s", prettyJSON.String())
				}
			}
		})
	}

	store.dropSchema(t)
	if err = store.client.Close(); err != nil {
		t.Logf("can't close store connection: %+v", err)
	}

}

type testCase struct {
	name           string
	method         string
	token          string
	ip             string
	cookie         *http.Cookie
	path           string
	prepareRequest func(req *http.Request)
	prepareDB      func(t *testing.T, stores *store)
	afterTest      func(t *testing.T, resp *http.Response)
	checkDB        func(t *testing.T, stores *store)
	expectedBody   string
	expectedStatus int
}

type store struct {
	client         *database.Client
	companyStorage dataprovider.CompaniesStorage
}

func (s *store) dropSchema(t *testing.T) {
	tx := s.client.MustBegin()
	tx.MustExec("DROP SCHEMA IF EXISTS " + s.client.SchemaName + " CASCADE;")
	err := tx.Commit()
	assert.NoError(t, err)
}

func configureIpCheckerMock(mock *service.IpCheckerMock) *service.IpCheckerMock {
	mock = mock.GetUserLocationMock.Set(func(_ context.Context, ip string) (location string, err error) {
		switch ip {
		case cyLocation:
			return "CY", err
		case usLocation:
			return "US", nil
		case bsLocation:
			return "", errors.New("unexpected response from external service")
		default:
			return "", ierr.UnknownLocation
		}
	})
	return mock
}

func configureMqMock(mock *service.MessageQueueMock) *service.MessageQueueMock {
	mock.NotifyCompanyUpdatedMock.Return(nil)
	return mock
}

func prepareRequest(body interface{}, location string) func(*http.Request) {
	return func(r *http.Request) {
		var reader io.Reader
		r.Header.Set("Content-type", "application/json")

		if location != "" {
			r.Header.Set("X-Real-Ip", location)
		}

		switch v := body.(type) {
		case string:
			reader = strings.NewReader(v)
		default:
			bt, err := json.Marshal(v)
			if err != nil {
				panic("can't marshal data into JSON: " + err.Error())
			}
			reader = bytes.NewReader(bt)
		}
		r.Body = ioutil.NopCloser(reader)
	}
}
