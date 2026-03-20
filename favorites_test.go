package pcloud

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListFavorites(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getfavourites" {
			http.Error(w, "wrong endpoint", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(favoritesResponse{
			Items: []Metadata{
				{Name: "fav1.jpg", FileID: 10},
				{Name: "fav2.pdf", FileID: 20},
			},
		})
	})

	items, err := c.ListFavorites(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("want 2 favorites, got %d", len(items))
	}
}

func TestAddFavorite(t *testing.T) {
	t.Parallel()
	var gotFileID string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotFileID = r.URL.Query().Get("fileid")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	if err := c.AddFavorite(t.Context(), 55); err != nil {
		t.Fatal(err)
	}
	if gotFileID != "55" {
		t.Fatalf("want fileid=55, got %q", gotFileID)
	}
}

func TestAddFavoriteByPath(t *testing.T) {
	t.Parallel()
	var gotPath string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Query().Get("path")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	if err := c.AddFavoriteByPath(t.Context(), "/docs/fav.pdf"); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/docs/fav.pdf" {
		t.Fatalf("want path=/docs/fav.pdf, got %q", gotPath)
	}
}

func TestRemoveFavorite(t *testing.T) {
	t.Parallel()
	var gotFileID string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotFileID = r.URL.Query().Get("fileid")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	if err := c.RemoveFavorite(t.Context(), 77); err != nil {
		t.Fatal(err)
	}
	if gotFileID != "77" {
		t.Fatalf("want fileid=77, got %q", gotFileID)
	}
}

func TestRemoveFavoriteByPath(t *testing.T) {
	t.Parallel()
	var gotPath string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Query().Get("path")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	if err := c.RemoveFavoriteByPath(t.Context(), "/docs/unfav.pdf"); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/docs/unfav.pdf" {
		t.Fatalf("want path=/docs/unfav.pdf, got %q", gotPath)
	}
}
