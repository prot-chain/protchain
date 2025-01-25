package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"protchain/internal/config"
	"protchain/internal/dep"
	"protchain/internal/value"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Handler func(w http.ResponseWriter, r *http.Request) *ServerResponse

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := h(w, r)
	responseByte, err := json.Marshal(response)
	if err != nil {
		writeErrorResponse(w, err, value.Error, "unable to marshal server response")
		return
	}
	writeJSONResponse(w, responseByte, response.StatusCode)
}

type API struct {
	Server *http.Server
	Config *config.Config
	Deps   *dep.Dependencies
}

// Serve starts the core service
func (a *API) Serve() error {
	a.Server = &http.Server{
		Addr:           fmt.Sprintf(":%d", a.Config.HttpPort),
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   40 * time.Second,
		Handler:        a.setUpServerHandler(),
		MaxHeaderBytes: 1024 * 1024,
	}

	slog.Info(fmt.Sprintf("Serving at port %d", a.Config.HttpPort))
	return a.Server.ListenAndServe()
}

func (a *API) Shutdown() error {
	return a.Server.Shutdown(context.Background())
}

// setUpServerHandler sets up handlers for the service
func (a *API) setUpServerHandler() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost, http.MethodGet, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions,
		},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", value.HeaderRequestID, value.HeaderRequestSource},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.WriteHeader(http.StatusNoContent)
	})
	mux.Mount("/auth", a.AuthRoutes(mux))
	mux.Mount("/protein", a.ProteinRoutes(mux))

	return mux
}
