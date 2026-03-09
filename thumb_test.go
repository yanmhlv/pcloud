package pcloud

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetThumbnail(t *testing.T) {
	var gotSize, gotFileID string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotSize = r.URL.Query().Get("size")
		gotFileID = r.URL.Query().Get("fileid")
		json.NewEncoder(w).Encode(FileLink{
			Path:  "/thumb/img.jpg",
			Hosts: []string{"cdn.pcloud.com"},
		})
	})

	link, err := c.GetThumbnail(t.Context(), 10, 200, 150, nil)
	if err != nil {
		t.Fatal(err)
	}
	if link.URL() == "" {
		t.Fatal("expected non-empty URL")
	}
	if gotSize != "200x150" {
		t.Fatalf("want size=200x150, got %q", gotSize)
	}
	if gotFileID != "10" {
		t.Fatalf("want fileid=10, got %q", gotFileID)
	}
}

func TestGetThumbnailByPath(t *testing.T) {
	var gotPath string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Query().Get("path")
		json.NewEncoder(w).Encode(FileLink{
			Path:  "/thumb/img.jpg",
			Hosts: []string{"cdn.pcloud.com"},
		})
	})

	_, err := c.GetThumbnailByPath(t.Context(), "/photo.jpg", 100, 100, nil)
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/photo.jpg" {
		t.Fatalf("want path=/photo.jpg, got %q", gotPath)
	}
}

func TestGetThumbnailWithOpts(t *testing.T) {
	var gotCrop, gotType string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotCrop = r.URL.Query().Get("crop")
		gotType = r.URL.Query().Get("type")
		json.NewEncoder(w).Encode(FileLink{
			Path:  "/thumb/img.png",
			Hosts: []string{"cdn.pcloud.com"},
		})
	})

	c.GetThumbnail(t.Context(), 1, 64, 64, &ThumbOpts{Crop: true, Type: "png"})
	if gotCrop != "1" {
		t.Fatalf("want crop=1, got %q", gotCrop)
	}
	if gotType != "png" {
		t.Fatalf("want type=png, got %q", gotType)
	}
}

func TestGetThumbnailInvalidSize(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(FileLink{})
	})

	if _, err := c.GetThumbnail(t.Context(), 1, 0, 100, nil); err == nil {
		t.Fatal("expected error for width=0")
	}
	if _, err := c.GetThumbnail(t.Context(), 1, 100, 2049, nil); err == nil {
		t.Fatal("expected error for height=2049")
	}
}
