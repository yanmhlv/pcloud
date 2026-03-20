package pcloud

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetFileLink(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("fileid") != "10" {
			http.Error(w, "bad fileid", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(FileLink{
			Path:  "/dl/file.txt",
			Hosts: []string{"cdn.pcloud.com"},
		})
	})

	link, err := c.GetFileLink(t.Context(), 10, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := "https://cdn.pcloud.com/dl/file.txt"
	if link.URL() != want {
		t.Fatalf("want %q, got %q", want, link.URL())
	}
}

func TestGetFileLinkByPath(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("path") != "/docs/report.pdf" {
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(FileLink{
			Path:  "/dl/report.pdf",
			Hosts: []string{"cdn.pcloud.com"},
		})
	})

	link, err := c.GetFileLinkByPath(t.Context(), "/docs/report.pdf", nil)
	if err != nil {
		t.Fatal(err)
	}
	if link.URL() == "" {
		t.Fatal("expected non-empty URL")
	}
}

func TestGetFileLinkOpts(t *testing.T) {
	t.Parallel()
	var gotForceDownload string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotForceDownload = r.URL.Query().Get("forcedownload")
		json.NewEncoder(w).Encode(FileLink{
			Path:  "/dl/file.txt",
			Hosts: []string{"cdn.pcloud.com"},
		})
	})

	c.GetFileLink(t.Context(), 10, &FileLinkOpts{ForceDownload: true})
	if gotForceDownload != "1" {
		t.Fatalf("want forcedownload=1, got %q", gotForceDownload)
	}
}

func TestGetFileLinkAPIError(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Error{Result: 2005, Message: "not found"})
	})

	_, err := c.GetFileLink(t.Context(), 999, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
