package api

import (
	ierr "github.com/IakimenkoD/xm-companies-service/internal/errors"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/dataprovider"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func getURLParam(r *http.Request, field string) (string, error) {
	param := chi.URLParam(r, field)
	if param == "" {
		return "", errors.New("empty param")
	}

	return param, nil
}

func getURLInt64(r *http.Request, field string) (int64, error) {
	param, err := getURLParam(r, field)
	if err != nil {
		return 0, err
	}

	int64Param, err := strconv.ParseInt(param, 10, 64)
	if err != nil || int64Param == 0 {
		return 0, ierr.InvalidParam
	}

	return int64Param, nil
}

func getQueryInt64Slice(r *http.Request, field string) ([]int64, error) {
	q := r.URL.Query()
	params := q[field]

	if len(params) == 0 {
		return nil, nil
	}

	var vals []int64

	for _, p := range params {
		slice := strings.Split(p, ",")
		if vals == nil {
			vals = make([]int64, 0, len(slice))
		}

		for _, s := range slice {
			if s == "" {
				continue
			}
			val, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, ierr.InvalidParam
			}
			vals = append(vals, val)
		}
	}

	return vals, nil
}

func getQueryStringSlice(r *http.Request, field string) ([]string, error) {
	q := r.URL.Query()
	params := q[field]

	if len(params) == 0 {
		return nil, nil
	}

	var vals []string

	for _, p := range params {
		slice := strings.Split(p, ",")
		if vals == nil {
			vals = make([]string, 0, len(slice))
		}

		for _, s := range slice {
			if s == "" {
				continue
			}
			vals = append(vals, s)
		}
	}

	return vals, nil
}

func parseCompaniesFilter(r *http.Request) (*dataprovider.CompanyFilter, error) {
	ids, err := getQueryInt64Slice(r, "ids")
	if err != nil {
		return nil, err
	}

	names, err := getQueryStringSlice(r, "names")
	if err != nil {
		return nil, err
	}

	codes, err := getQueryStringSlice(r, "codes")
	if err != nil {
		return nil, err
	}

	countries, err := getQueryStringSlice(r, "countries")
	if err != nil {
		return nil, err
	}

	websites, err := getQueryStringSlice(r, "websites")
	if err != nil {
		return nil, err
	}

	phones, err := getQueryStringSlice(r, "phones")
	if err != nil {
		return nil, err
	}

	for n := range phones {
		phones[n] = normalizePhoneNumber(phones[n])
	}

	return dataprovider.NewCompanyFilter().
		ByIDs(ids...).
		ByNames(names...).
		ByCodes(codes...).
		ByCountries(toLowerCase(countries)...).
		ByWebsites(toLowerCase(websites)...).
		ByPhones(phones...), nil
}

func normalizePhoneNumber(phoneNumber string) string {
	phoneNumber = strings.TrimSpace(phoneNumber)

	if len(phoneNumber) == 0 {
		return ""
	}

	if !strings.HasPrefix(phoneNumber, "+") {
		phoneNumber = "+" + phoneNumber
	}

	return phoneNumber
}

func respondError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ierr.CompanyNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case errors.Is(err, ierr.InvalidParam), errors.Is(err, ierr.WrongRequest), errors.Is(err, io.EOF), errors.Is(err, ierr.CompanyExists):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func toLowerCase(values []string) []string {
	result := make([]string, 0, len(values))
	for _, v := range values {
		result = append(values, strings.ToLower(v))
	}
	return result
}
