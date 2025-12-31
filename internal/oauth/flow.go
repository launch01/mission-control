package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/launch01/mission-control/internal/config"
	"github.com/launch01/mission-control/internal/logging"
	"github.com/launch01/mission-control/internal/storage"
)

// TokenResponse represents the OAuth token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// AuthFlow handles the OAuth 2.0 authorization flow
type AuthFlow struct {
	cfg           *config.Config
	storage       *storage.TokenStorage
	callbackChan  chan string
	server        *http.Server
}

// NewAuthFlow creates a new OAuth flow handler
func NewAuthFlow(cfg *config.Config) (*AuthFlow, error) {
	store, err := storage.NewTokenStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to create token storage: %w", err)
	}

	return &AuthFlow{
		cfg:          cfg,
		storage:      store,
		callbackChan: make(chan string, 1),
	}, nil
}

// Login initiates the OAuth login flow
func (f *AuthFlow) Login(ctx context.Context) error {
	// Generate PKCE parameters
	verifier, err := GenerateCodeVerifier()
	if err != nil {
		return fmt.Errorf("failed to generate code verifier: %w", err)
	}
	challenge := GenerateCodeChallenge(verifier)

	state, err := GenerateState()
	if err != nil {
		return fmt.Errorf("failed to generate state: %w", err)
	}

	// Build authorization URL
	authURL := f.buildAuthURL(challenge, state)
	logging.Info("Opening browser for authorization...")
	logging.Info("Please visit: %s", authURL)

	// Start callback server
	port := f.getCallbackPort()
	if err := f.startCallbackServer(port, state); err != nil {
		return fmt.Errorf("failed to start callback server: %w", err)
	}
	defer f.stopCallbackServer()

	// Open browser
	if err := openBrowser(authURL); err != nil {
		logging.Error("Failed to open browser automatically: %v", err)
		logging.Info("Please open the URL manually in your browser")
	}

	// Wait for callback
	select {
	case code := <-f.callbackChan:
		logging.Info("Authorization code received")
		return f.exchangeCodeForToken(code, verifier)
	case <-ctx.Done():
		return fmt.Errorf("login canceled")
	case <-time.After(5 * time.Minute):
		return fmt.Errorf("login timeout")
	}
}

// RefreshToken refreshes the access token using the refresh token
func (f *AuthFlow) RefreshToken(ctx context.Context) error {
	token, err := f.storage.LoadToken()
	if err != nil {
		return fmt.Errorf("failed to load token: %w", err)
	}

	if token.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", f.cfg.HubSpot.ClientID)
	data.Set("refresh_token", token.RefreshToken)

	if f.cfg.HubSpot.ClientSecret != "" {
		data.Set("client_secret", f.cfg.HubSpot.ClientSecret)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", config.HubSpotTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	// Update stored token
	token.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		token.RefreshToken = tokenResp.RefreshToken
	}
	token.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	if err := f.storage.SaveToken(token); err != nil {
		return fmt.Errorf("failed to save refreshed token: %w", err)
	}

	logging.Info("Token refreshed successfully")
	return nil
}

func (f *AuthFlow) buildAuthURL(challenge, state string) string {
	params := url.Values{}
	params.Set("client_id", f.cfg.HubSpot.ClientID)
	params.Set("redirect_uri", f.cfg.HubSpot.RedirectURI)
	params.Set("scope", f.cfg.HubSpot.Scopes)
	params.Set("state", state)
	params.Set("code_challenge", challenge)
	params.Set("code_challenge_method", "S256")

	return config.HubSpotAuthURL + "?" + params.Encode()
}

func (f *AuthFlow) startCallbackServer(port, expectedState string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")
		errorParam := r.URL.Query().Get("error")

		if errorParam != "" {
			http.Error(w, "Authorization failed: "+errorParam, http.StatusBadRequest)
			close(f.callbackChan)
			return
		}

		if state != expectedState {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			close(f.callbackChan)
			return
		}

		if code == "" {
			http.Error(w, "No authorization code received", http.StatusBadRequest)
			close(f.callbackChan)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<html>
			<body>
				<h1>Authorization Successful!</h1>
				<p>You can close this window and return to the terminal.</p>
			</body>
			</html>
		`)

		f.callbackChan <- code
	})

	f.server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		if err := f.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("Callback server error: %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (f *AuthFlow) stopCallbackServer() {
	if f.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		f.server.Shutdown(ctx)
	}
}

func (f *AuthFlow) exchangeCodeForToken(code, verifier string) error {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", f.cfg.HubSpot.ClientID)
	data.Set("redirect_uri", f.cfg.HubSpot.RedirectURI)
	data.Set("code", code)
	data.Set("code_verifier", verifier)

	if f.cfg.HubSpot.ClientSecret != "" {
		data.Set("client_secret", f.cfg.HubSpot.ClientSecret)
	}

	req, err := http.NewRequest("POST", config.HubSpotTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	// Save token
	token := &storage.Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}

	if err := f.storage.SaveToken(token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	logging.Info("Authentication successful!")
	return nil
}

func (f *AuthFlow) getCallbackPort() string {
	u, err := url.Parse(f.cfg.HubSpot.RedirectURI)
	if err != nil {
		return config.DefaultCallbackPort
	}
	port := u.Port()
	if port == "" {
		return config.DefaultCallbackPort
	}
	return port
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch {
	case isCommandAvailable("xdg-open"):
		cmd = "xdg-open"
		args = []string{url}
	case isCommandAvailable("open"):
		cmd = "open"
		args = []string{url}
	case isCommandAvailable("start"):
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return fmt.Errorf("unable to open browser automatically")
	}

	return exec.Command(cmd, args...).Start()
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
