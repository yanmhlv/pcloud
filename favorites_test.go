package pcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListFavorites(t *testing.T) {
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

	items, err := c.ListFavorites(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("want 2 favorites, got %d", len(items))
	}
}

func TestAddFavorite(t *testing.T) {
	var gotFileID string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotFileID = r.URL.Query().Get("fileid")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	if err := c.AddFavorite(context.Background(), 55); err != nil {
		t.Fatal(err)
	}
	if gotFileID != "55" {
		t.Fatalf("want fileid=55, got %q", gotFileID)
	}
}

func TestAddFavoriteByPath(t *testing.T) {
	var gotPath string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Query().Get("path")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	if err := c.AddFavoriteByPath(context.Background(), "/docs/fav.pdf"); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/docs/fav.pdf" {
		t.Fatalf("want path=/docs/fav.pdf, got %q", gotPath)
	}
}

func TestRemoveFavorite(t *testing.T) {
	var gotFileID string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotFileID = r.URL.Query().Get("fileid")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	if err := c.RemoveFavorite(context.Background(), 77); err != nil {
		t.Fatal(err)
	}
	if gotFileID != "77" {
		t.Fatalf("want fileid=77, got %q", gotFileID)
	}
}

func TestRemoveFavoriteByPath(t *testing.T) {
	var gotPath string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Query().Get("path")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	if err := c.RemoveFavoriteByPath(context.Background(), "/docs/unfav.pdf"); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/docs/unfav.pdf" {
		t.Fatalf("want path=/docs/unfav.pdf, got %q", gotPath)
	}
}
