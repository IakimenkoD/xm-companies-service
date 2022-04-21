package api

import (
	"encoding/json"
	ierr "github.com/IakimenkoD/xm-companies-service/internal/errors"
	"github.com/IakimenkoD/xm-companies-service/internal/model"
	"github.com/IakimenkoD/xm-companies-service/internal/repository/dataprovider"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func (srv *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) getCompanies(w http.ResponseWriter, r *http.Request) {
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

	filter := dataprovider.NewCompanyFilter().ByIDs(id)
	company, err := srv.controller.GetCompanies(ctx, filter)
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

	company := &model.Company{}
	if err := json.NewDecoder(r.Body).Decode(company); err != nil {
		respondError(w, errors.Wrap(ierr.WrongRequest, err.Error()))
		return
	}
	if err := company.CheckFields(); err != nil {
		respondError(w, err)
		return
	}

	id, err := srv.controller.CreateCompany(ctx, company)
	if err != nil {
		respondError(w, err)
		return
	}
	company.ID = id

	w.Header().Set("Location", "/company/"+strconv.FormatInt(id, 10))
	w.WriteHeader(http.StatusCreated)
}

func (srv *Server) updateCompany(w http.ResponseWriter, r *http.Request) {
	id, err := getURLInt64(r, "companyID")
	if err != nil {
		respondError(w, err)
		return
	}

	ctx := r.Context()
	data := &model.Company{}
	if err = json.NewDecoder(r.Body).Decode(data); err != nil {
		respondError(w, err)
		return
	}
	data.ID = id

	if err = srv.controller.UpdateCompany(ctx, data); err != nil {
		respondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (srv *Server) patchCompany(w http.ResponseWriter, r *http.Request) {
	id, err := getURLInt64(r, "companyID")
	if err != nil {
		respondError(w, err)
		return
	}

	ctx := r.Context()
	company := &model.Company{}
	if err = json.NewDecoder(r.Body).Decode(company); err != nil {
		respondError(w, err)
		return
	}
	company.ID = id

	company, err = srv.controller.PatchCompany(ctx, company)
	if err != nil {
		respondError(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(company); err != nil {
		respondError(w, err)
		return
	}
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
