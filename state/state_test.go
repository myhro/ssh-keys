package state

import (
	"os"
	"path/filepath"
	"testing"
)

const perm = 0o600

func TestReadMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "authorized_keys.ssh-keys")

	data, err := New(path).Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.ETag != "" || data.SHA256 != "" {
		t.Fatalf("expected an empty state, got %+v", data)
	}
}

func TestReadUnreadableFile(t *testing.T) {
	_, err := New(t.TempDir()).Read()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestReadInvalidContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "not json",
			content: `"abc123"` + "\n",
		},
		{
			name:    "truncated",
			content: `{"etag": "abc123`,
		},
		{
			name:    "wrong types",
			content: `{"etag": 1, "sha256": []}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "authorized_keys.ssh-keys")
			err := os.WriteFile(path, []byte(tc.content), perm)
			if err != nil {
				t.Fatalf("writing file: %v", err)
			}

			data, err := New(path).Read()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if data.ETag != "" || data.SHA256 != "" {
				t.Fatalf("expected an empty state, got %+v", data)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "authorized_keys.ssh-keys")
	file := New(path)

	err := file.Write(&Data{ETag: `W/"abc123"`, SHA256: Digest([]byte("keys\n"))})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := file.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.ETag != `W/"abc123"` {
		t.Fatalf(`expected W/"abc123", got %q`, data.ETag)
	}
	if data.SHA256 != Digest([]byte("keys\n")) {
		t.Fatalf("expected the recorded digest, got %q", data.SHA256)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("checking file: %v", err)
	}
	if info.Mode().Perm() != perm {
		t.Fatalf("expected mode %o, got %o", perm, info.Mode().Perm())
	}

	_, err = os.Stat(path + ".tmp")
	if !os.IsNotExist(err) {
		t.Fatal("expected the temporary file to be gone")
	}
}

func TestWriteMissingDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "authorized_keys.ssh-keys")

	err := New(path).Write(&Data{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestDigest(t *testing.T) {
	const want = "71e10a1ea66438a092639ff7b67f115a049d2a4684287f401d9c2bf0e51f8eab"

	got := Digest([]byte("keys\n"))
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
	if Digest([]byte("keys")) == got {
		t.Fatal("expected trailing whitespace to change the digest")
	}
}
