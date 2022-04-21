package api

import (
	mw "github.com/IakimenkoD/xm-companies-service/internal/api/middleware"
	"github.com/IakimenkoD/xm-companies-service/internal/config"
	"github.com/IakimenkoD/xm-companies-service/internal/controller"
	"github.com/IakimenkoD/xm-companies-service/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

type Server struct {
	*http.Server
	controller controller.CompaniesService
	ipChecker  service.IpChecker
	cfg        *config.Config
}

func NewServer(
	cfg *config.Config,
	controller controller.CompaniesService,
	ipChecker service.IpChecker,

) (*Server, error) {
	srv := &Server{
		Server: &http.Server{
			Addr:         cfg.API.Address,
			ReadTimeout:  cfg.API.ReadTimeout,
			WriteTimeout: cfg.API.WriteTimeout,
		},
		cfg:        cfg,
		controller: controller,
		ipChecker:  ipChecker,
	}

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Route("/internal", func(r chi.Router) {
		r.Post("/signin", srv.signIn)
		r.Get("/health", srv.health)

	})

	r.Route("/api/v1/companies", func(r chi.Router) {
		r.Get("/", srv.getCompanies)
		r.Route("/{companyID}", func(r chi.Router) {
			r.Get("/", srv.getCompanyByID)
			r.Put("/", srv.updateCompany)
			r.Patch("/", srv.patchCompany)
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.CheckIPAddress(srv.ipChecker))
			r.Use(mw.CheckAuth(srv.cfg.API.JWTKey))

			r.Post("/", srv.createCompany)
			r.Delete("/{companyID}", srv.deleteCompany)
		})
	})

	srv.Handler = r

	return srv, nil
}
