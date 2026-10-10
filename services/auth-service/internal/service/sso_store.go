package service

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"

	"com.ecommerce/auth-service/internal/model"
)

type stateEntry struct {
	redirectURL string
	expiresAt   time.Time
}

type tokenEntry struct {
	tokens    model.KeycloakTokenResponse
	expiresAt time.Time
}

type SsoSessionStore struct {
	mu        sync.RWMutex
	states    map[string]stateEntry
	tickets   map[string]tokenEntry
	stateTTL  time.Duration
	ticketTTL time.Duration
}

func NewSsoSessionStore() *SsoSessionStore {
	store := &SsoSessionStore{
		states:    make(map[string]stateEntry),
		tickets:   make(map[string]tokenEntry),
		stateTTL:  5 * time.Minute,
		ticketTTL: 2 * time.Minute,
	}

	// Purge routine
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			store.purge()
		}
	}()

	return store
}

func (s *SsoSessionStore) purge() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for k, v := range s.states {
		if v.expiresAt.Before(now) {
			delete(s.states, k)
		}
	}
	for k, v := range s.tickets {
		if v.expiresAt.Before(now) {
			delete(s.tickets, k)
		}
	}
}

func randomToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (s *SsoSessionStore) CreateLoginState(frontendRedirect string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := randomToken()
	s.states[state] = stateEntry{
		redirectURL: frontendRedirect,
		expiresAt:   time.Now().Add(s.stateTTL),
	}
	return state
}

func (s *SsoSessionStore) ConsumeLoginState(state string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.states[state]
	if !ok {
		return "", false
	}
	delete(s.states, state)
	if entry.expiresAt.Before(time.Now()) {
		return "", false
	}
	return entry.redirectURL, true
}

func (s *SsoSessionStore) StoreTokens(tokens model.KeycloakTokenResponse) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticket := randomToken()
	s.tickets[ticket] = tokenEntry{
		tokens:    tokens,
		expiresAt: time.Now().Add(s.ticketTTL),
	}
	return ticket
}

func (s *SsoSessionStore) ConsumeTokens(ticket string) (*model.KeycloakTokenResponse, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.tickets[ticket]
	if !ok {
		return nil, false
	}
	delete(s.tickets, ticket)
	if entry.expiresAt.Before(time.Now()) {
		return nil, false
	}
	return &entry.tokens, true
}
