package routers

import (
	"net/http"

	"github.com/diaku/eht-demo/internal/handlers"
	"github.com/diaku/eht-demo/internal/middlewares"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

type Deps struct {
	Auth *handlers.AuthHanlder
	Eth  *handlers.EthHandler

	AuthMiddleware func(http.Handler) http.Handler
}

func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	// global middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Logger)

	// Auth
	r.Route("/auth", func(r chi.Router) {
		r.Get("/nonce", d.Auth.Nouce)
		r.Post("/verify", d.Auth.VerifySIWE)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(d.AuthMiddleware)

		r.Route("/eth", func(r chi.Router) {
			r.Get("/chain", d.Eth.Chain)
			r.Get("/balance", d.Eth.Balance)
		})

		// Example "me" endpoint later (reads wallet address from context)
		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			addr, _ := r.Context().Value(middlewares.CtxWalletAddress).(string)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"address":"` + addr + `"}`))
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	return r
}
