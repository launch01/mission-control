package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "mission-control"
	tokenKey    = "hubspot-token"
)

// Token represents an OAuth token
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// TokenStorage handles secure storage of OAuth tokens
type TokenStorage struct {
	useKeyring bool
	filePath   string
}

// NewTokenStorage creates a new token storage
func NewTokenStorage() (*TokenStorage, error) {
	// Try to use keyring first
	if isKeyringAvailable() {
		return &TokenStorage{useKeyring: true}, nil
	}

	// Fall back to file storage with warning
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "mission-control")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	filePath := filepath.Join(configDir, "token.json")
	return &TokenStorage{
		useKeyring: false,
		filePath:   filePath,
	}, nil
}

// SaveToken saves the token securely
func (s *TokenStorage) SaveToken(token *Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	if s.useKeyring {
		return keyring.Set(serviceName, tokenKey, string(data))
	}

	// File storage with restricted permissions
	if err := os.WriteFile(s.filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// LoadToken loads the token from secure storage
func (s *TokenStorage) LoadToken() (*Token, error) {
	var data string
	var err error

	if s.useKeyring {
		data, err = keyring.Get(serviceName, tokenKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get token from keyring: %w", err)
		}
	} else {
		bytes, err := os.ReadFile(s.filePath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("no token found - please run 'mission-control auth login' first")
			}
			return nil, fmt.Errorf("failed to read token file: %w", err)
		}
		data = string(bytes)
	}

	var token Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// DeleteToken deletes the stored token
func (s *TokenStorage) DeleteToken() error {
	if s.useKeyring {
		return keyring.Delete(serviceName, tokenKey)
	}

	if err := os.Remove(s.filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete token file: %w", err)
	}

	return nil
}

// IsExpired checks if the token is expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsExpiringSoon checks if the token will expire within the given duration
func (t *Token) IsExpiringSoon(d time.Duration) bool {
	return time.Now().Add(d).After(t.ExpiresAt)
}

func isKeyringAvailable() bool {
	// Test if keyring is available by trying to set/get/delete a test value
	testKey := "test-availability"
	err := keyring.Set(serviceName, testKey, "test")
	if err != nil {
		return false
	}
	keyring.Delete(serviceName, testKey)
	return true
}
