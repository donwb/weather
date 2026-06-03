package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// tokenData is what we persist to disk between requests/restarts.
type tokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"` // unix seconds
}

// tokenStore is a concurrency-safe, file-backed holder for the Netatmo
// OAuth tokens. Netatmo rotates the refresh token on every use, so the
// rotated value must be persisted or auth is lost on the next call.
type tokenStore struct {
	mu   sync.Mutex
	path string
	data tokenData
}

var tokens *tokenStore

func newTokenStore(path string) *tokenStore {
	ts := &tokenStore{path: path}
	ts.load()
	return ts
}

func (ts *tokenStore) load() {
	b, err := os.ReadFile(ts.path)
	if err != nil {
		return // no file yet; caller bootstraps from env
	}
	var d tokenData
	if err := json.Unmarshal(b, &d); err != nil {
		logError(err, "parsing token file")
		return
	}
	ts.data = d
}

// save writes atomically so a crash mid-write can't corrupt the token file.
func (ts *tokenStore) save() {
	b, err := json.MarshalIndent(ts.data, "", "  ")
	if err != nil {
		logError(err, "marshalling tokens")
		return
	}
	tmp := ts.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		logError(err, "writing token file")
		return
	}
	if err := os.Rename(tmp, ts.path); err != nil {
		logError(err, "replacing token file")
	}
}

// bootstrap seeds the store from env vars, but only if we don't already
// have a refresh token on disk (the disk copy is newer after any refresh).
func (ts *tokenStore) bootstrap(access, refresh string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.data.RefreshToken != "" || refresh == "" {
		return
	}
	ts.data.AccessToken = access
	ts.data.RefreshToken = refresh
	ts.data.ExpiresAt = 0 // force a refresh on first use
	ts.save()
}

func (ts *tokenStore) set(access, refresh string, expiresIn int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.data.AccessToken = access
	if refresh != "" {
		ts.data.RefreshToken = refresh
	}
	if expiresIn > 0 {
		ts.data.ExpiresAt = time.Now().Unix() + int64(expiresIn)
	}
	ts.save()
}

func (ts *tokenStore) snapshot() tokenData {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.data
}

// expired reports whether the access token is missing or within 60s of expiry.
func (ts *tokenStore) expired() bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.data.AccessToken == "" {
		return true
	}
	return time.Now().Unix() >= ts.data.ExpiresAt-60
}
