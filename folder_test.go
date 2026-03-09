package pcloud

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestListFolder(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("folderid") != "0" {
			http.Error(w, "bad folderid", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(metadataResponse{
			Metadata: Metadata{Name: "root", IsFolder: true},
		})
	})

	meta, err := c.ListFolder(t.Context(), 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Name != "root" {
		t.Fatalf("want root, got %q", meta.Name)
	}
}

func TestListFolderByPath(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("path") != "/myfolder" {
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(metadataResponse{
			Metadata: Metadata{Name: "myfolder", IsFolder: true},
		})
	})

	meta, err := c.ListFolderByPath(t.Context(), "/myfolder", nil)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Name != "myfolder" {
		t.Fatalf("want myfolder, got %q", meta.Name)
	}
}

func TestListFolderAPIError(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(metadataResponse{Error: Error{Result: 2005, Message: "not found"}})
	})

	_, err := c.ListFolder(t.Context(), 999, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListFolderOpts(t *testing.T) {
	var gotRecursive string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotRecursive = r.URL.Query().Get("recursive")
		json.NewEncoder(w).Encode(metadataResponse{})
	})

	c.ListFolder(t.Context(), 0, &ListFolderOpts{Recursive: true})
	if gotRecursive != "1" {
		t.Fatalf("want recursive=1, got %q", gotRecursive)
	}
}

func TestStatFolder(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(metadataResponse{
			Metadata: Metadata{Name: "docs", IsFolder: true, FolderID: 42},
		})
	})

	meta, err := c.StatFolder(t.Context(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if meta.FolderID != 42 {
		t.Fatalf("want folderid=42, got %d", meta.FolderID)
	}
}

func TestStatFolderByPath(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("path") != "/docs" {
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(metadataResponse{
			Metadata: Metadata{Name: "docs", IsFolder: true},
		})
	})

	meta, err := c.StatFolderByPath(t.Context(), "/docs")
	if err != nil {
		t.Fatal(err)
	}
	if meta.Name != "docs" {
		t.Fatalf("want docs, got %q", meta.Name)
	}
}

func TestCreateFolder(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(metadataResponse{
			Metadata: Metadata{Name: "new-dir", IsFolder: true, FolderID: 10},
		})
	})

	meta, err := c.CreateFolder(t.Context(), 0, "new-dir")
	if err != nil {
		t.Fatal(err)
	}
	if meta.FolderID != 10 {
		t.Fatalf("want folderid=10, got %d", meta.FolderID)
	}
}

func TestWalkYieldsContents(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(metadataResponse{
			Metadata: Metadata{
				IsFolder: true,
				Contents: []Metadata{
					{Name: "file1.txt", IsFolder: false},
					{Name: "sub", IsFolder: true, Contents: []Metadata{
						{Name: "file2.txt", IsFolder: false},
					}},
				},
			},
		})
	})

	var names []string
	for item, err := range c.Walk(t.Context(), 0) {
		if err != nil {
			t.Fatal(err)
		}
		names = append(names, item.Name)
	}

	if len(names) != 3 {
		t.Fatalf("want 3 items, got %d: %v", len(names), names)
	}
}

func TestWalkEarlyBreak(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(metadataResponse{
			Metadata: Metadata{
				IsFolder: true,
				Contents: []Metadata{
					{Name: "a"},
					{Name: "b"},
					{Name: "c"},
				},
			},
		})
	})

	var count int
	for _, err := range c.Walk(t.Context(), 0) {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if count == 1 {
			break
		}
	}
	if count != 1 {
		t.Fatalf("want 1, got %d", count)
	}
}
