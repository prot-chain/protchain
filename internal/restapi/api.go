package restapi

import (
	"context"
	"encoding/json"
	"fmt"
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
		Addr:           fmt.Sprintf(":%d", a.Config.Port),
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   40 * time.Second,
		Handler:        a.setUpServerHandler(),
		MaxHeaderBytes: 1024 * 1024,
	}

	return a.Server.ListenAndServe()
}

func (a *API) Shutdown() error {
	return a.Server.Shutdown(context.Background())
}

// setUpServerHandler sets up handlers for the service
func (a *API) setUpServerHandler() http.Handler {
	mux := chi.NewRouter()

	mux.Route("/", func(r chi.Router) {
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(60 * time.Second))
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPatch, http.MethodPut, http.MethodDelete},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-For", value.HeaderRequestID, value.HeaderRequestSource},
			AllowCredentials: true,
			MaxAge:           300,
		}))

		r.Use(RequestTracing)

		r.Mount("/mono", a.MonoRoutes())
		r.Mount("/termii", a.TermiiRoutes())
		r.Mount("/twilio", a.TwilioRoutes())
		r.Mount("/lidya", a.LidyaRoutes())
		r.Mount("/crc", a.CrcRoutes())
		r.Mount("/cellulant", a.CellulantRoutes())
		r.Mount("/filestack", a.FileStackRoutes())
		r.Mount("/identitypass", a.IdentitpPassRoutes())
		r.Mount("/paystack", a.PayStackRoutes())
		r.Mount("/cyberpay", a.CyberPayRoutes())
		r.Mount("/remita", a.RemitaRoutes())
		r.Mount("/calendy", a.CalendyRoutes())
		r.Mount("/mbs", a.MbsRoutes())
		r.Mount("/okra", a.OkraRoutes())
		r.Mount("/elastic", a.ElasticRoutes())
		r.Mount("/blacklist", a.BlackListRoutes())
		r.Mount("/flutterwave", a.FlutterwaveRoutes())
		r.Mount("/whatsapp", a.WhatsAppRoutes())
		r.Mount("/bitnob", a.BitnobRoutes())
		r.Mount("/providus", a.ProvidusRoutes())
		r.Mount("/graph", a.GraphCmsRoutes())
		r.Mount("/providus_transfer", a.ProvidusTransferRoutes())
		r.Mount("/stripe", a.StripeRoutes())
		r.Mount("/periculum", a.PericulumRoutes())
		r.Mount("/payment", a.PaymentRoutes())
		r.Mount("/nlng", a.NLNGRoutes())
		r.Mount("/iprocess", a.IprocessRoutes())
		r.Mount("/rubicon", a.RubiconRoutes())
		// go a.SmsConsumer()
	})

	mux.Route("/webhook", func(webhookRouter chi.Router) {
		webhookRouter.Use(middleware.RealIP)
		webhookRouter.Use(middleware.Logger)
		webhookRouter.Use(middleware.Recoverer)
		webhookRouter.Use(middleware.Timeout(60 * time.Second))
		webhookRouter.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPatch, http.MethodPut, http.MethodDelete},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-For", value.HeaderRequestID, value.HeaderRequestSource},
			AllowCredentials: true,
			MaxAge:           300,
		}))

		webhookRouter.Use(RequestTracingII)

		webhookRouter.Mount("/payment", a.PaymentWebHookRoutes())
	})
	return mux
}
