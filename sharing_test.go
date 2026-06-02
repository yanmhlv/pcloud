package pcloud

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestShareFolder(t *testing.T) {
	t.Parallel()
	var gotEmail string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotEmail = r.URL.Query().Get("mail")
		json.NewEncoder(w).Encode(Share{ShareID: 9})
	})

	share, err := c.ShareFolder(t.Context(), 1, "bob@example.com", SharePermissions{CanRead: true}, nil)
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
	t.Parallel()
	var gotPath string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Query().Get("path")
		json.NewEncoder(w).Encode(Share{ShareID: 11})
	})

	_, err := c.ShareFolderByPath(t.Context(), "/shared", "alice@example.com", SharePermissions{CanRead: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/shared" {
		t.Fatalf("want path=/shared, got %q", gotPath)
	}
}

func TestListShares(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(listSharesResponse{
			Shares:   []Share{{ShareID: 1}},
			Requests: []Share{{ShareID: 2}},
		})
	})

	shares, err := c.ListShares(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(shares.Active) != 1 || shares.Active[0].ShareID != 1 {
		t.Fatalf("unexpected active shares: %v", shares.Active)
	}
	if len(shares.Requests) != 1 || shares.Requests[0].ShareID != 2 {
		t.Fatalf("unexpected requests: %v", shares.Requests)
	}
}
