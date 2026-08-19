package fetch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"abc123"`)
		w.Write([]byte("keys\n"))
	}))
	defer server.Close()

	res, err := New().Get(context.Background(), server.URL, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.NotModified {
		t.Fatal("expected a modified response")
	}
	if string(res.Body) != "keys\n" {
		t.Fatalf("expected %q, got %q", "keys\n", string(res.Body))
	}
	if res.ETag != `"abc123"` {
		t.Fatalf(`expected "abc123", got %q`, res.ETag)
	}
}

func TestGetNotModified(t *testing.T) {
	sent := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sent = r.Header.Get("If-None-Match")
		w.WriteHeader(http.StatusNotModified)
	}))
	defer server.Close()

	res, err := New().Get(context.Background(), server.URL, `"abc123"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.NotModified {
		t.Fatal("expected a not modified response")
	}
	if len(res.Body) != 0 {
		t.Fatalf("expected an empty body, got %q", string(res.Body))
	}
	if sent != `"abc123"` {
		t.Fatalf(`expected the If-None-Match header to be "abc123", got %q`, sent)
	}
}

func TestGetWithoutETag(t *testing.T) {
	header := "unset"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, ok := r.Header["If-None-Match"]
		if ok {
			header = value[0]
		} else {
			header = ""
		}
		w.Write([]byte("keys\n"))
	}))
	defer server.Close()

	_, err := New().Get(context.Background(), server.URL, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if header != "" {
		t.Fatalf("expected no If-None-Match header, got %q", header)
	}
}

func TestGetErrors(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		invalid bool
	}{
		{
			name:   "not found",
			status: http.StatusNotFound,
			body:   "404 page not found\n",
		},
		{
			name:   "server error",
			status: http.StatusInternalServerError,
		},
		{
			name:   "empty body",
			status: http.StatusOK,
		},
		{
			name:    "unreachable server",
			invalid: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.body))
			}))
			defer server.Close()

			url := server.URL
			if tc.invalid {
				server.Close()
			}

			res, err := New().Get(context.Background(), url, "")
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if res != nil {
				t.Fatalf("expected no response, got %+v", res)
			}
		})
	}
}
