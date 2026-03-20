package pcloud

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestCreateFilePublicLink(t *testing.T) {
	t.Parallel()
	var gotFileID string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotFileID = r.URL.Query().Get("fileid")
		json.NewEncoder(w).Encode(PublicLink{LinkID: 7, Code: "abc"})
	})

	link, err := c.CreateFilePublicLink(t.Context(), 42, nil)
	if err != nil {
		t.Fatal(err)
	}
	if link.LinkID != 7 {
		t.Fatalf("want linkid=7, got %d", link.LinkID)
	}
	if gotFileID != "42" {
		t.Fatalf("want fileid=42, got %q", gotFileID)
	}
}

func TestCreateFilePublicLinkByPath(t *testing.T) {
	t.Parallel()
	var gotPath string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Query().Get("path")
		json.NewEncoder(w).Encode(PublicLink{LinkID: 8})
	})

	_, err := c.CreateFilePublicLinkByPath(t.Context(), "/report.pdf", nil)
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/report.pdf" {
		t.Fatalf("want path=/report.pdf, got %q", gotPath)
	}
}

func TestCreateFilePublicLinkOpts(t *testing.T) {
	t.Parallel()
	var gotExpire, gotMaxDownloads string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotExpire = r.URL.Query().Get("expire")
		gotMaxDownloads = r.URL.Query().Get("maxdownloads")
		json.NewEncoder(w).Encode(PublicLink{})
	})

	c.CreateFilePublicLink(t.Context(), 1, &PublicLinkOpts{
		ExpireAt:     time.Unix(1700000000, 0),
		MaxDownloads: 5,
	})
	if gotExpire != "1700000000" {
		t.Fatalf("want expire=1700000000, got %q", gotExpire)
	}
	if gotMaxDownloads != "5" {
		t.Fatalf("want maxdownloads=5, got %q", gotMaxDownloads)
	}
}

func TestListPublicLinks(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(listPublicLinksResponse{
			PubLinks: []PublicLink{{LinkID: 1}, {LinkID: 2}},
		})
	})

	links, err := c.ListPublicLinks(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 2 {
		t.Fatalf("want 2 links, got %d", len(links))
	}
}

func TestDeletePublicLink(t *testing.T) {
	t.Parallel()
	var gotLinkID string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotLinkID = r.URL.Query().Get("linkid")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	if err := c.DeletePublicLink(t.Context(), 15); err != nil {
		t.Fatal(err)
	}
	if gotLinkID != "15" {
		t.Fatalf("want linkid=15, got %q", gotLinkID)
	}
}
