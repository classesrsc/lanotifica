package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	t.Parallel()

	secret := "test-secret-token"
	nextCalled := false

	next := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}

	handler := AuthMiddleware(secret, next)

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+secret)

	rr := httptest.NewRecorder()
	handler(rr, req)

	if !nextCalled {
		t.Error("Expected next handler to be called")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	t.Parallel()

	secret := "test-secret-token"
	nextCalled := false

	next := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}

	handler := AuthMiddleware(secret, next)

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	// No Authorization header

	rr := httptest.NewRecorder()
	handler(rr, req)

	if nextCalled {
		t.Error("Expected next handler NOT to be called")
	}
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	t.Parallel()

	secret := "test-secret-token"
	nextCalled := false

	next := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}

	handler := AuthMiddleware(secret, next)

	testCases := []struct {
		name   string
		header string
	}{
		{"Basic auth", "Basic " + secret},
		{"No space", "Bearer" + secret},
		{"Only token", secret},
		{"Empty bearer", "Bearer"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
			req.Header.Set("Authorization", tc.header)

			rr := httptest.NewRecorder()
			handler(rr, req)

			if nextCalled {
				t.Error("Expected next handler NOT to be called")
			}
			if rr.Code != http.StatusUnauthorized {
				t.Errorf("Expected status 401, got %d", rr.Code)
			}
		})
	}
}

func TestAuthMiddleware_WrongToken(t *testing.T) {
	t.Parallel()

	secret := "correct-secret"
	nextCalled := false

	next := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}

	handler := AuthMiddleware(secret, next)

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("Authorization", "Bearer wrong-token")

	rr := httptest.NewRecorder()
	handler(rr, req)

	if nextCalled {
		t.Error("Expected next handler NOT to be called")
	}
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_EmptyToken(t *testing.T) {
	t.Parallel()

	secret := "test-secret"
	nextCalled := false

	next := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}

	handler := AuthMiddleware(secret, next)

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("Authorization", "Bearer ")

	rr := httptest.NewRecorder()
	handler(rr, req)

	if nextCalled {
		t.Error("Expected next handler NOT to be called")
	}
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}
