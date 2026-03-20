package pcloud

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestSearch(t *testing.T) {
	t.Parallel()
	var gotQuery string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("query")
		json.NewEncoder(w).Encode(searchResponse{
			Items: []Metadata{
				{Name: "report.pdf", FileID: 5},
				{Name: "report-final.pdf", FileID: 6},
			},
		})
	})

	items, err := c.Search(t.Context(), "report", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("want 2 results, got %d", len(items))
	}
	if gotQuery != "report" {
		t.Fatalf("want query=report, got %q", gotQuery)
	}
}

func TestSearchWithOpts(t *testing.T) {
	t.Parallel()
	var gotFolderID, gotRecursive string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotFolderID = r.URL.Query().Get("folderid")
		gotRecursive = r.URL.Query().Get("recursive")
		json.NewEncoder(w).Encode(searchResponse{})
	})

	c.Search(t.Context(), "img", &SearchOpts{FolderID: 10, Recursive: true})
	if gotFolderID != "10" {
		t.Fatalf("want folderid=10, got %q", gotFolderID)
	}
	if gotRecursive != "1" {
		t.Fatalf("want recursive=1, got %q", gotRecursive)
	}
}

func TestSearchAPIError(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(searchResponse{Error: Error{Result: 2003, Message: "access denied"}})
	})

	_, err := c.Search(t.Context(), "secret", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
