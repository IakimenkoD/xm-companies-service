package api

import (
	"encoding/json"
	"github.com/IakimenkoD/xm-companies-service/internal/model"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func (srv *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) getCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	filter, err := parseCompaniesFilter(r)
	if err != nil {
		respondError(w, err)
		return
	}

	companies, err := srv.controller.GetCompanies(ctx, filter)
	if err != nil {
		respondError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(companies); err != nil {
		respondError(w, err)
		return
	}
}

func (srv *Server) getCompanyByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := getURLInt64(r, "companyID")
	if err != nil {
		respondError(w, err)
		return
	}

	company, err := srv.controller.GetCompanyByID(ctx, id)
	if err != nil {
		respondError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(company); err != nil {
		respondError(w, err)
		return
	}
}

func (srv *Server) createCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := &model.Company{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		respondError(w, errors.Wrap(err, "decoding request"))
		return
	}

	id, err := srv.controller.CreateCompany(ctx, data)
	if err != nil {
		respondError(w, err)
		return
	}
	w.Header().Set("Location", "/company/"+strconv.FormatInt(id, 10))
	w.WriteHeader(http.StatusCreated)
}

func (srv *Server) updateCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := &model.Company{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		respondError(w, err)
		return
	}
	if err := srv.controller.UpdateCompany(ctx, data); err != nil {
		respondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (srv *Server) deleteCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := getURLInt64(r, "companyID")
	if err != nil {
		respondError(w, err)
		return
	}

	if err = srv.controller.DeleteCompany(ctx, id); err != nil {
		respondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
