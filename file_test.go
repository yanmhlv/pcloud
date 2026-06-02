package pcloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
)

func TestStat(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("fileid") != "123" {
			http.Error(w, "bad fileid", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(metadataResponse{
			Metadata: Metadata{Name: "photo.jpg", FileID: 123},
		})
	})

	meta, err := c.Stat(t.Context(), 123)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Name != "photo.jpg" {
		t.Fatalf("want photo.jpg, got %q", meta.Name)
	}
}

func TestStatByPath(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("path") != "/photo.jpg" {
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(metadataResponse{
			Metadata: Metadata{Name: "photo.jpg"},
		})
	})

	meta, err := c.StatByPath(t.Context(), "/photo.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if meta.Name != "photo.jpg" {
		t.Fatalf("want photo.jpg, got %q", meta.Name)
	}
}

func TestDeleteFile(t *testing.T) {
	t.Parallel()
	var gotFileID string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotFileID = r.URL.Query().Get("fileid")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	if err := c.DeleteFile(t.Context(), 55); err != nil {
		t.Fatal(err)
	}
	if gotFileID != "55" {
		t.Fatalf("want fileid=55, got %q", gotFileID)
	}
}

func TestDeleteFileByPath(t *testing.T) {
	t.Parallel()
	var gotPath string
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Query().Get("path")
		json.NewEncoder(w).Encode(Error{Result: 0})
	})

	if err := c.DeleteFileByPath(t.Context(), "/old.txt"); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/old.txt" {
		t.Fatalf("want path=/old.txt, got %q", gotPath)
	}
}

func TestUpload(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(uploadResponse{
			Metadata: []Metadata{{Name: "hello.txt", FileID: 77}},
		})
	})

	content := bytes.NewReader([]byte("hello"))
	meta, err := c.Upload(t.Context(), 0, "hello.txt", content, nil)
	if err != nil {
		t.Fatal(err)
	}
	if meta.FileID != 77 {
		t.Fatalf("want fileid=77, got %d", meta.FileID)
	}
}

func TestProgressReaderCumulative(t *testing.T) {
	t.Parallel()
	var calls []int64
	pr := &progressReader{
		rc:    io.NopCloser(strings.NewReader("0123456789")),
		total: 10,
		onProgress: func(transferred, total int64) {
			calls = append(calls, transferred)
		},
	}

	buf := make([]byte, 5)
	if _, err := pr.Read(buf); err != nil {
		t.Fatal(err)
	}
	if _, err := pr.Read(buf); err != nil && !errors.Is(err, io.EOF) {
		t.Fatal(err)
	}

	if len(calls) < 2 {
		t.Fatalf("expected 2 progress callbacks, got %d", len(calls))
	}
	if calls[0] != 5 {
		t.Fatalf("first callback: want 5, got %d", calls[0])
	}
	if calls[1] != 10 {
		t.Fatalf("second callback: want 10, got %d", calls[1])
	}
}

func TestConcurrentDownloadAndSetHTTPClient(t *testing.T) {
	t.Parallel()
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(fileLinkResponse{
			FileLink: FileLink{Path: "/file", Hosts: []string{"127.0.0.1"}},
		})
	})

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			rc, err := c.Download(t.Context(), 1, nil)
			if err == nil {
				rc.Close()
			}
		})
		wg.Go(func() {
			rc, err := c.DownloadByPath(t.Context(), "/file", nil)
			if err == nil {
				rc.Close()
			}
		})
		wg.Go(func() {
			c.SetHTTPClient(&http.Client{})
		})
	}
	wg.Wait()
}
