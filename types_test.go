package pcloud

import (
	"encoding/json"
	"testing"
)

func TestHashUnmarshalString(t *testing.T) {
	t.Parallel()
	var h Hash
	if err := json.Unmarshal([]byte(`"abc123"`), &h); err != nil {
		t.Fatal(err)
	}
	if h != "abc123" {
		t.Fatalf("want abc123, got %q", h)
	}
}

func TestHashUnmarshalNumber(t *testing.T) {
	t.Parallel()
	var h Hash
	if err := json.Unmarshal([]byte(`12345678`), &h); err != nil {
		t.Fatal(err)
	}
	if h != "12345678" {
		t.Fatalf("want 12345678, got %q", h)
	}
}

func TestHashUnmarshalInvalid(t *testing.T) {
	t.Parallel()
	var h Hash
	if err := json.Unmarshal([]byte(`true`), &h); err == nil {
		t.Fatal("expected error for boolean hash")
	}
}

func TestTimeUnmarshalValid(t *testing.T) {
	t.Parallel()
	var ts Time
	if err := json.Unmarshal([]byte(`"Mon, 02 Jan 2006 15:04:05 -0700"`), &ts); err != nil {
		t.Fatal(err)
	}
	if ts.IsZero() {
		t.Fatal("expected non-zero time")
	}
}

func TestTimeUnmarshalEmpty(t *testing.T) {
	t.Parallel()
	var ts Time
	if err := json.Unmarshal([]byte(`""`), &ts); err != nil {
		t.Fatal(err)
	}
	if !ts.IsZero() {
		t.Fatal("expected zero time for empty string")
	}
}

func TestErrorErr(t *testing.T) {
	t.Parallel()
	e := &Error{Result: 0}
	if e.Err() != nil {
		t.Fatal("result=0 should return nil error")
	}

	e = &Error{Result: ErrNotFound, Message: "not found"}
	if e.Err() == nil {
		t.Fatal("result!=0 should return error")
	}
	if e.Error() != "pcloud error 2005: not found" {
		t.Fatalf("unexpected error string: %s", e.Error())
	}
}

func TestErrorErrNoMessage(t *testing.T) {
	t.Parallel()
	e := &Error{Result: ErrAuthRequired}
	if e.Error() != "pcloud error 1000" {
		t.Fatalf("unexpected error string: %s", e.Error())
	}
}

func TestFileLinkURL(t *testing.T) {
	t.Parallel()
	f := &FileLink{Path: "/dl/file.txt", Hosts: []string{"cdn1.pcloud.com"}}
	want := "https://cdn1.pcloud.com/dl/file.txt"
	if f.URL() != want {
		t.Fatalf("want %q, got %q", want, f.URL())
	}
}

func TestFileLinkURLNoHosts(t *testing.T) {
	t.Parallel()
	f := &FileLink{Path: "/dl/file.txt"}
	if f.URL() != "" {
		t.Fatalf("expected empty URL, got %q", f.URL())
	}
}
