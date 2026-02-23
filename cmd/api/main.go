package main

import (
	"context"
	"log"
	"net/http"

	"github.com/diaku/eht-demo/config"
	"github.com/diaku/eht-demo/internal/handlers"
	"github.com/diaku/eht-demo/internal/middlewares"
	"github.com/diaku/eht-demo/internal/routers"
	"github.com/diaku/eht-demo/internal/usecase"
	"github.com/diaku/eht-demo/pkg/utils"
	"github.com/joho/godotenv"
)

// TEMP: stubs so the app compiles while we implement real logic
type stubAuthUC struct{}

func (s stubAuthUC) NewNonce(_ context.Context) (string, error) { return "stub-nonce", nil } // fix signature when you implement
func (s stubAuthUC) VerifySIWE(_ context.Context, _ string, _ string) (string, string, error) {
	return "stub-token", "0x0", nil
}

type stubEthUC struct{}

func (s stubEthUC) ChainInfo(_ context.Context) (string, error)                 { return "11155111", nil }
func (s stubEthUC) Balance(_ context.Context, _ string) (string, string, error) { return "0", "0", nil }

// TEMP: token verifier stub
type stubVerifier struct{}

func (s stubVerifier) Verify(_ string) (string, error) { return "0x0", nil }

func main() {
	// loading env
	_ = godotenv.Load()

	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	nonceStore := usecase.NewNonceStore()
	jwtSvc := utils.NewJWTService(cfg.JWTSecret, cfg.JWTExpiry, "go-eth-demo")

	// Usecases
	authUC := usecase.NewAuthUsecase(usecase.AuthDeps{
		NonceStore: nonceStore,
		JWT:        jwtSvc,
		Domain:     cfg.SiweDomain,
		URI:        cfg.SiweURI,
		ChainID:    cfg.ChainID,
		TimeSkew:   cfg.SiweTimeSkew,
		NonceTTL:   cfg.SiweNonceTTL,
	})

	ethUC, err := usecase.NewEthUsecase(cfg.AnkrRPCURL)
	if err != nil {
		log.Fatal(err)
	}

	// Handlers
	authHandler := handlers.NewAuthHandler(authUC)
	ethHandler := handlers.NewEthHandler(ethUC)

	verifier := stubVerifier{}
	authMW := middlewares.RequireAuth(verifier)

	r := routers.NewRouter(routers.Deps{
		Auth:           authHandler,
		Eth:            ethHandler,
		AuthMiddleware: authMW,
	})

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
