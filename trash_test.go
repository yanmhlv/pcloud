package pcloud

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListTrash(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/trash_list" {
			http.Error(w, "wrong endpoint", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(trashListResponse{
			Items: []TrashItem{
				{Name: "old.txt", FileID: 1},
				{Name: "archive", IsFolder: true, FolderID: 2},
			},
		})
	})

	items, err := c.ListTrash(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("want 2 items, got %d", len(items))
	}
	if items[0].Name != "old.txt" {
		t.Fatalf("want old.txt, got %q", items[0].Name)
	}
}

func TestRestoreFromTrash(t *testing.T) {
	var gotFileID string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotFileID = r.URL.Query().Get("fileid")
		json.NewEncoder(w).Encode(metadataResponse{
			Metadata: Metadata{Name: "old.txt"},
		})
	})

	meta, err := c.RestoreFromTrash(t.Context(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Name != "old.txt" {
		t.Fatalf("want old.txt, got %q", meta.Name)
	}
	if gotFileID != "42" {
		t.Fatalf("want fileid=42, got %q", gotFileID)
	}
}

func TestRestoreFromTrashByPath(t *testing.T) {
	var gotPath string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Query().Get("path")
		json.NewEncoder(w).Encode(metadataResponse{
			Metadata: Metadata{Name: "old.txt"},
		})
	})

	_, err := c.RestoreFromTrashByPath(t.Context(), "/trash/old.txt")
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/trash/old.txt" {
		t.Fatalf("want path=/trash/old.txt, got %q", gotPath)
	}
}

func TestEmptyTrash(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/trash_clear" {
			http.Error(w, "wrong endpoint", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	if err := c.EmptyTrash(t.Context()); err != nil {
		t.Fatal(err)
	}
}
