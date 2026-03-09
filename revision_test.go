package pcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListRevisions(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("fileid") != "20" {
			http.Error(w, "bad fileid", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(revisionsResponse{
			Revisions: []Revision{
				{RevisionID: 1, Size: 100},
				{RevisionID: 2, Size: 200},
			},
		})
	})

	revs, err := c.ListRevisions(context.Background(), 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(revs) != 2 {
		t.Fatalf("want 2 revisions, got %d", len(revs))
	}
}

func TestListRevisionsByPath(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("path") != "/file.txt" {
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(revisionsResponse{
			Revisions: []Revision{{RevisionID: 5, Size: 512}},
		})
	})

	revs, err := c.ListRevisionsByPath(context.Background(), "/file.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(revs) != 1 || revs[0].RevisionID != 5 {
		t.Fatalf("unexpected revisions: %v", revs)
	}
}

func TestRevertRevision(t *testing.T) {
	var gotRevisionID string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotRevisionID = r.URL.Query().Get("revisionid")
		json.NewEncoder(w).Encode(metadataResponse{
			Metadata: Metadata{Name: "file.txt"},
		})
	})

	meta, err := c.RevertRevision(context.Background(), 20, 3)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Name != "file.txt" {
		t.Fatalf("want file.txt, got %q", meta.Name)
	}
	if gotRevisionID != "3" {
		t.Fatalf("want revisionid=3, got %q", gotRevisionID)
	}
}
