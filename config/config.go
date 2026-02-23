package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Server
	Port string `envconfig:"PORT" default:"8080"`

	// Ethereum / RPC
	AnkrRPCURL string `envconfig:"ANKR_RPC_URL" required:"true"`
	ChainID    int64  `envconfig:"CHAIN_ID" default:"11155111"` // Sepolia by default

	// SIWE settings
	SiweDomain    string        `envconfig:"SIWE_DOMAIN" required:"true"`                   // e.g. "localhost" or "myapp.com"
	SiweStatement string        `envconfig:"SIWE_STATEMENT" default:"Sign in to GoEthDemo"` // human-readable statement
	SiweURI       string        `envconfig:"SIWE_URI" default:"http://localhost:8080"`      // must match the SIWE message URI
	SiweNonceTTL  time.Duration `envconfig:"SIWE_NONCE_TTL" default:"5m"`                   // nonce expiry
	SiweTimeSkew  time.Duration `envconfig:"SIWE_TIME_SKEW" default:"2m"`                   // allow small client/server clock skew

	// Auth / JWT
	JWTSecret     string        `envconfig:"JWT_SECRET" required:"true"`
	JWTExpiry     time.Duration `envconfig:"JWT_EXPIRY" default:"24h"`
	AllowedOrigin string        `envconfig:"ALLOWED_ORIGIN" default:"*"` // for CORS in demo
}

func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, fmt.Errorf("environment config failed to process, error: %w", err)
	}
	return &config, nil
}
