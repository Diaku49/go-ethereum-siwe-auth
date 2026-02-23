package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/diaku/eht-demo/pkg/utils"
	"github.com/spruceid/siwe-go"
)

type AuthUC interface {
	NewNonce(ctx context.Context) (string, error)
	VerifySIWE(ctx context.Context, message string, signature string) (token string, walletAddress string, err error)
}

type AuthUsecase struct {
	nonceStore *NonceStore
	jwt        *utils.JWTService

	// SIWE constraints (from config)
	domain   string
	uri      string
	chainID  int64
	timeSkew time.Duration
	nonceTTL time.Duration
}

type AuthDeps struct {
	NonceStore *NonceStore
	JWT        *utils.JWTService

	Domain   string
	URI      string
	ChainID  int64
	TimeSkew time.Duration
	NonceTTL time.Duration
}

func NewAuthUsecase(deps AuthDeps) *AuthUsecase {
	return &AuthUsecase{
		nonceStore: deps.NonceStore,
		jwt:        deps.JWT,
		domain:     deps.Domain,
		uri:        deps.URI,
		chainID:    deps.ChainID,
		timeSkew:   deps.TimeSkew,
		nonceTTL:   deps.NonceTTL,
	}
}

func (a *AuthUsecase) NewNonce(ctx context.Context) (string, error) {
	_ = ctx
	return a.nonceStore.New(a.nonceTTL)
}

func (a *AuthUsecase) VerifySIWE(ctx context.Context, messageStr string, signature string) (token string, walletAddress string, err error) {
	_ = ctx

	msg, err := siwe.ParseMessage(messageStr)
	if err != nil {
		return "", "", fmt.Errorf("invalid siwe message: %w", err)
	}

	// 1) Check the message is intended for THIS app
	if msg.GetDomain() != a.domain {
		return "", "", fmt.Errorf("siwe domain mismatch")
	}
	u := msg.GetURI()
	if u.String() != a.uri {
		return "", "", fmt.Errorf("siwe uri mismatch: got=%q want=%q", u.String(), a.uri)
	}
	if a.chainID > 0 && int64(msg.GetChainID()) != a.chainID {
		return "", "", fmt.Errorf("siwe chainId mismatch")
	}

	// 2) Nonce must exist, not expired, not used
	nonce := msg.GetNonce()
	if nonce == "" {
		return "", "", fmt.Errorf("missing nonce")
	}
	if !a.nonceStore.Valid(nonce) {
		return "", "", fmt.Errorf("invalid or expired nonce")
	}

	// 3) Time validity (+ small skew allowance)
	now := time.Now()
	ok, err := msg.ValidAt(now.UTC())
	if err != nil || !ok {
		return "", "", fmt.Errorf("siwe message time constraints invalid")
	}

	// 4) Cryptographic verification (EIP-191)
	// Verify(signature, optionalDomain, optionalNonce, optionalTimestamp)
	d := a.domain
	n := nonce
	t := now
	if _, err := msg.Verify(signature, &d, &n, &t); err != nil {
		return "", "", fmt.Errorf("signature verification failed")
	}

	// 5) Consume nonce (single-use)
	if !a.nonceStore.Consume(nonce) {
		return "", "", fmt.Errorf("nonce already used")
	}

	// 6) Mint JWT
	addr := msg.GetAddress().String()
	jwtToken, err := a.jwt.Issue(addr)
	if err != nil {
		return "", "", fmt.Errorf("failed to issue token: %w", err)
	}

	return jwtToken, addr, nil
}
