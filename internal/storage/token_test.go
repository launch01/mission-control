package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTokenExpiry(t *testing.T) {
	token := &Token{
		AccessToken:  "test-access",
		RefreshToken: "test-refresh",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	if token.IsExpired() {
		t.Error("Token should not be expired")
	}

	if token.IsExpiringSoon(2 * time.Hour) {
		t.Error("Token should not be expiring soon for 2 hours")
	}

	if !token.IsExpiringSoon(30 * time.Minute) {
		t.Error("Token should be expiring soon for 30 minutes")
	}

	expiredToken := &Token{
		AccessToken:  "test-access",
		RefreshToken: "test-refresh",
		ExpiresAt:    time.Now().Add(-1 * time.Hour),
	}

	if !expiredToken.IsExpired() {
		t.Error("Token should be expired")
	}
}

func TestFileStorage(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	storage := &TokenStorage{
		useKeyring: false,
		filePath:   filepath.Join(tmpDir, "token.json"),
	}

	token := &Token{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	// Test save
	err := storage.SaveToken(token)
	if err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	// Test load
	loaded, err := storage.LoadToken()
	if err != nil {
		t.Fatalf("LoadToken() error = %v", err)
	}

	if loaded.AccessToken != token.AccessToken {
		t.Errorf("AccessToken = %v, want %v", loaded.AccessToken, token.AccessToken)
	}

	if loaded.RefreshToken != token.RefreshToken {
		t.Errorf("RefreshToken = %v, want %v", loaded.RefreshToken, token.RefreshToken)
	}

	// Test delete
	err = storage.DeleteToken()
	if err != nil {
		t.Fatalf("DeleteToken() error = %v", err)
	}

	// Should not exist
	_, err = storage.LoadToken()
	if err == nil {
		t.Error("LoadToken() should return error after deletion")
	}
}

func TestFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()

	storage := &TokenStorage{
		useKeyring: false,
		filePath:   filepath.Join(tmpDir, "token.json"),
	}

	token := &Token{
		AccessToken:  "test-access",
		RefreshToken: "test-refresh",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	err := storage.SaveToken(token)
	if err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	// Check file permissions
	info, err := os.Stat(storage.filePath)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	mode := info.Mode().Perm()
	if mode != 0600 {
		t.Errorf("File permissions = %o, want 0600", mode)
	}
}
