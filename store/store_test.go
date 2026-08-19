package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "authorized_keys")

	file := New(path)
	data, err := file.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != nil {
		t.Fatalf("expected no data, got %q", data)
	}

	err = os.WriteFile(path, []byte(" keys \n"), perm)
	if err != nil {
		t.Fatalf("writing file: %v", err)
	}

	data, err = file.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != " keys \n" {
		t.Fatalf("expected the bytes to be returned as-is, got %q", data)
	}
}

func TestReadUnreadableFile(t *testing.T) {
	path := t.TempDir()

	_, err := New(path).Read()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "authorized_keys")
	file := New(path)

	err := file.Write([]byte("first\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = file.Write([]byte("second\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}
	if string(data) != "second\n" {
		t.Fatalf("expected %q, got %q", "second\n", string(data))
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

func TestWriteReplacesStaleTemporaryFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "authorized_keys")
	tmp := path + ".tmp"

	err := os.WriteFile(tmp, []byte("leftovers\n"), perm)
	if err != nil {
		t.Fatalf("writing file: %v", err)
	}

	err = New(path).Write([]byte("keys\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}
	if string(data) != "keys\n" {
		t.Fatalf("expected %q, got %q", "keys\n", string(data))
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("checking file: %v", err)
	}
	if info.Mode().Perm() != perm {
		t.Fatalf("expected mode %o, got %o", perm, info.Mode().Perm())
	}
}

func TestWriteMissingDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "authorized_keys")

	err := New(path).Write([]byte("keys\n"))
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
