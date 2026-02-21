package pcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
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
	if err := c.SetRateLimit(DefaultRPM - 1); err == nil {
		t.Fatal("expected error for rate below DefaultRPM")
	}
}

func TestSetRateLimitValid(t *testing.T) {
	c := NewClient("")
	if err := c.SetRateLimit(DefaultRPM); err != nil {
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
	if err := c.do(context.Background(), "ping", url.Values{}, &resp); err != nil {
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
	err := c.do(context.Background(), "stat", url.Values{}, &resp)
	if err == nil {
		t.Fatal("expected error from API result=2005")
	}
	if err.Error() != "pcloud error 2005: not found" {
		t.Fatalf("unexpected error: %s", err)
	}
}
