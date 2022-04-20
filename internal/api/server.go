package api

import (
	"github.com/IakimenkoD/xm-companies-service/internal/config"
	"github.com/IakimenkoD/xm-companies-service/internal/controller"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

type Server struct {
	*http.Server
	controller controller.CompaniesService
	cfg        *config.Config
}

func NewServer(
	cfg *config.Config,
	controller controller.CompaniesService,
) (*Server, error) {
	srv := &Server{
		Server: &http.Server{
			Addr:         cfg.API.Address,
			ReadTimeout:  cfg.API.ReadTimeout,
			WriteTimeout: cfg.API.WriteTimeout,
		},
		cfg:        cfg,
		controller: controller,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Route("/internal", func(r chi.Router) {
		r.Get("/health", srv.health)

	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/company", func(r chi.Router) {
				r.Get("/", srv.getCompany)
				r.Post("/", srv.createCompany)
				r.Get("/{companyID}", srv.getCompanyByID)
				r.Put("/{companyID}", srv.updateCompany)
				r.Delete("/{companyID}", srv.deleteCompany)
			})
		})
	})

	srv.Handler = r

	return srv, nil
}
