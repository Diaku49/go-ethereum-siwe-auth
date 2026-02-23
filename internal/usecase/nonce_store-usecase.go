package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

type NonceStore struct {
	mu sync.Mutex
	m  map[string]nonceEntry
}

type nonceEntry struct {
	expiresAt time.Time
	used      bool
}

func NewNonceStore() *NonceStore {
	return &NonceStore{
		m: make(map[string]nonceEntry),
	}
}

func (s *NonceStore) New(ttl time.Duration) (string, error) {
	nonce, err := generateNonce(18)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[nonce] = nonceEntry{
		expiresAt: time.Now().Add(ttl),
		used:      false,
	}
	return nonce, nil
}

func (s *NonceStore) Valid(nonce string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	ent, ok := s.m[nonce]
	if !ok {
		return false
	}
	if ent.used {
		return false
	}
	if time.Now().After(ent.expiresAt) {
		delete(s.m, nonce)
		return false
	}
	return true
}

func (s *NonceStore) Consume(nonce string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	ent, ok := s.m[nonce]
	if !ok {
		return false
	}
	if ent.used {
		return false
	}
	if time.Now().After(ent.expiresAt) {
		delete(s.m, nonce)
		return false
	}

	ent.used = true
	s.m[nonce] = ent
	return true
}

func generateNonce(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
