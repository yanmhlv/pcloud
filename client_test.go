package pcloud

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
)

func TestNewClientDefaultURL(t *testing.T) {
	c := NewClient("")
	if c.baseURL != BaseURLUS {
		t.Fatalf("want %q, got %q", BaseURLUS, c.baseURL)
	}
}

func TestSetRateLimitBelowMin(t *testing.T) {
	c := NewClient("")
	if err := c.SetRateLimit(MinRPM - 1); err == nil {
		t.Fatal("expected error for rate below MinRPM")
	}
}

func TestSetRateLimitValid(t *testing.T) {
	c := NewClient("")
	if err := c.SetRateLimit(MinRPM); err != nil {
		t.Fatal(err)
	}
}

func TestRequestAuthInjected(t *testing.T) {
	var gotAuth string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.URL.Query().Get("auth")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	var resp Error
	if err := c.do(t.Context(), "ping", url.Values{}, &resp); err != nil {
		t.Fatal(err)
	}
	if gotAuth != "test-token" {
		t.Fatalf("want auth=test-token, got %q", gotAuth)
	}
}

func TestRequestAPIError(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Error{Result: 2005, Message: "not found"})
	})

	var resp Error
	err := c.do(t.Context(), "stat", url.Values{}, &resp)
	if err == nil {
		t.Fatal("expected error from API result=2005")
	}
	if err.Error() != "pcloud error 2005: not found" {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestConcurrentLoginAndRequest(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(loginResponse{Auth: "tok"})
	})

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			c.Login(t.Context(), "u", "p")
		}()
		go func() {
			defer wg.Done()
			var resp Error
			c.do(t.Context(), "ping", url.Values{}, &resp)
		}()
	}
	wg.Wait()
}

func TestRequestHTTP500(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<html>Internal Server Error</html>"))
	})

	var resp Error
	err := c.do(t.Context(), "test", url.Values{}, &resp)
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Fatalf("error should mention status code: %s", err)
	}
}

func TestRequestHTTP200Valid(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	var resp Error
	if err := c.do(t.Context(), "test", url.Values{}, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
