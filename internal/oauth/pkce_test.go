package oauth

import (
	"testing"
)

func TestGenerateCodeVerifier(t *testing.T) {
	verifier, err := GenerateCodeVerifier()
	if err != nil {
		t.Fatalf("GenerateCodeVerifier() error = %v", err)
	}

	if len(verifier) == 0 {
		t.Error("GenerateCodeVerifier() returned empty string")
	}

	// Verify it's URL-safe base64
	if len(verifier) < 43 || len(verifier) > 128 {
		t.Errorf("GenerateCodeVerifier() length = %d, want between 43 and 128", len(verifier))
	}
}

func TestGenerateCodeChallenge(t *testing.T) {
	verifier := "test-verifier-12345"
	challenge := GenerateCodeChallenge(verifier)

	if len(challenge) == 0 {
		t.Error("GenerateCodeChallenge() returned empty string")
	}

	// Should be deterministic
	challenge2 := GenerateCodeChallenge(verifier)
	if challenge != challenge2 {
		t.Error("GenerateCodeChallenge() not deterministic")
	}
}

func TestGenerateState(t *testing.T) {
	state, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState() error = %v", err)
	}

	if len(state) == 0 {
		t.Error("GenerateState() returned empty string")
	}

	// Should be random
	state2, _ := GenerateState()
	if state == state2 {
		t.Error("GenerateState() not generating unique values")
	}
}
