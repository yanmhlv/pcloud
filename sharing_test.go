package pcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestShareFolder(t *testing.T) {
	var gotEmail string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotEmail = r.URL.Query().Get("mail")
		json.NewEncoder(w).Encode(Share{ShareID: 9})
	})

	share, err := c.ShareFolder(context.Background(), 1, "bob@example.com", SharePermissions{CanRead: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if share.ShareID != 9 {
		t.Fatalf("want shareid=9, got %d", share.ShareID)
	}
	if gotEmail != "bob@example.com" {
		t.Fatalf("want mail=bob@example.com, got %q", gotEmail)
	}
}

func TestShareFolderByPath(t *testing.T) {
	var gotPath string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Query().Get("path")
		json.NewEncoder(w).Encode(Share{ShareID: 11})
	})

	_, err := c.ShareFolderByPath(context.Background(), "/shared", "alice@example.com", SharePermissions{CanRead: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/shared" {
		t.Fatalf("want path=/shared, got %q", gotPath)
	}
}

func TestListShares(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(listSharesResponse{
			Shares:   []Share{{ShareID: 1}},
			Requests: []Share{{ShareID: 2}},
		})
	})

	shares, requests, err := c.ListShares(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(shares) != 1 || shares[0].ShareID != 1 {
		t.Fatalf("unexpected shares: %v", shares)
	}
	if len(requests) != 1 || requests[0].ShareID != 2 {
		t.Fatalf("unexpected requests: %v", requests)
	}
}
