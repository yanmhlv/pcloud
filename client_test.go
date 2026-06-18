package pcloud

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"

	"golang.org/x/oauth2"
)

func TestNewClientDefaultURL(t *testing.T) {
	t.Parallel()
	c := NewClient("")
	if c.baseURL != BaseURLUS {
		t.Fatalf("want %q, got %q", BaseURLUS, c.baseURL)
	}
}

func TestSetRateLimitBelowMin(t *testing.T) {
	t.Parallel()
	c := NewClient("")
	if err := c.SetRateLimit(MinRPM - 1); err == nil {
		t.Fatal("expected error for rate below MinRPM")
	}
}

func TestSetRateLimitValid(t *testing.T) {
	t.Parallel()
	c := NewClient("")
	if err := c.SetRateLimit(MinRPM); err != nil {
		t.Fatal(err)
	}
}

func TestRequestAuthInjected(t *testing.T) {
	t.Parallel()
	var gotAuth string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.URL.Query().Get("auth")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	var resp Error
	if err := c.doGet(t.Context(), "ping", url.Values{}, &resp); err != nil {
		t.Fatal(err)
	}
	if gotAuth != "test-token" {
		t.Fatalf("want auth=test-token, got %q", gotAuth)
	}
}

func TestRequestAPIError(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Error{Result: 2005, Message: "not found"})
	})

	var resp Error
	err := c.doGet(t.Context(), "stat", url.Values{}, &resp)
	if err == nil {
		t.Fatal("expected error from API result=2005")
	}
	if err.Error() != "pcloud error 2005: not found" {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestConcurrentLoginAndRequest(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(loginResponse{Auth: "tok"})
	})

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			c.Login(t.Context(), "u", "p")
		})
		wg.Go(func() {
			var resp Error
			c.doGet(t.Context(), "ping", url.Values{}, &resp)
		})
	}
	wg.Wait()
}

func TestRequestHTTP500(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<html>Internal Server Error</html>"))
	})

	var resp Error
	err := c.doGet(t.Context(), "test", url.Values{}, &resp)
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Fatalf("error should mention status code: %s", err)
	}
}

func TestRequestHTTP200Valid(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	var resp Error
	if err := c.doGet(t.Context(), "test", url.Values{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

type failingTokenSource struct{}

func (failingTokenSource) Token() (*oauth2.Token, error) {
	return nil, errors.New("token error")
}

func TestSetAuthWithTokenSource(t *testing.T) {
	t.Parallel()
	c := NewClient("")
	c.SetTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "oauth-token"}))

	params := url.Values{}
	if err := c.setAuth(params); err != nil {
		t.Fatal(err)
	}
	if params.Get("auth") != "oauth-token" {
		t.Fatalf("want auth=oauth-token, got %q", params.Get("auth"))
	}
}

func TestSetAuthWithFailingTokenSource(t *testing.T) {
	t.Parallel()
	c := NewClient("")
	c.SetTokenSource(failingTokenSource{})

	params := url.Values{}
	if err := c.setAuth(params); err == nil {
		t.Fatal("expected error from failing token source")
	}
}

func TestSetAuthTokenSourcePrecedence(t *testing.T) {
	t.Parallel()
	c := NewClient("")
	c.auth = "password-auth"
	c.SetTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "oauth-token"}))

	params := url.Values{}
	if err := c.setAuth(params); err != nil {
		t.Fatal(err)
	}
	if params.Get("auth") != "oauth-token" {
		t.Fatalf("want oauth-token (precedence), got %q", params.Get("auth"))
	}
}
