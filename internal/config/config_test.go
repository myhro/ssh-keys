package config

import (
	"errors"
	"os/user"
	"path/filepath"
	"testing"
)

type fakeUser struct {
	err      error
	username string
}

func (f fakeUser) Current() (*user.User, error) {
	return &user.User{Username: f.username}, f.err
}

func setenv(t *testing.T, vars map[string]string) {
	t.Helper()

	for _, name := range []string{"HOME", "SSH_KEYS_FILE", "SSH_KEYS_URL", "SSH_KEYS_USER"} {
		t.Setenv(name, vars[name])
	}
}

func TestLoadDefaults(t *testing.T) {
	home := t.TempDir()
	setenv(t, map[string]string{
		"HOME": home,
	})

	cfg := &Config{user: fakeUser{username: "myhro"}}
	err := cfg.load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	file := filepath.Join(home, ".ssh", "authorized_keys")
	if cfg.URL != "https://github.com/myhro.keys" {
		t.Fatalf("expected the default URL, got %q", cfg.URL)
	}
	if cfg.File != file {
		t.Fatalf("expected %q, got %q", file, cfg.File)
	}
	if cfg.ETagFile != file+".etag" {
		t.Fatalf("expected %q, got %q", file+".etag", cfg.ETagFile)
	}
}

func TestLoadOverrides(t *testing.T) {
	setenv(t, map[string]string{
		"HOME":          t.TempDir(),
		"SSH_KEYS_FILE": "/srv/authorized_keys",
		"SSH_KEYS_URL":  "https://example.com/keys.txt",
		"SSH_KEYS_USER": "ignored",
	})

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.URL != "https://example.com/keys.txt" {
		t.Fatalf("expected the URL override, got %q", cfg.URL)
	}
	if cfg.File != "/srv/authorized_keys" {
		t.Fatalf("expected the file override, got %q", cfg.File)
	}
	if cfg.ETagFile != "/srv/authorized_keys.etag" {
		t.Fatalf("expected the etag file next to the keys file, got %q", cfg.ETagFile)
	}
}

func TestLoadUserOverride(t *testing.T) {
	setenv(t, map[string]string{
		"HOME":          t.TempDir(),
		"SSH_KEYS_USER": "octocat",
	})

	cfg := &Config{user: fakeUser{username: "myhro"}}
	err := cfg.load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.URL != "https://github.com/octocat.keys" {
		t.Fatalf("expected the octocat URL, got %q", cfg.URL)
	}
}

func TestLoadUnknownUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		err      error
	}{
		{
			name: "lookup failure",
			err:  errors.New("user: lookup failed"),
		},
		{
			name: "empty username",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setenv(t, map[string]string{
				"HOME": t.TempDir(),
			})

			cfg := &Config{user: fakeUser{err: tc.err, username: tc.username}}
			err := cfg.load()
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	}
}

func TestLoadETagFileWithoutKeysFile(t *testing.T) {
	setenv(t, map[string]string{
		"HOME": t.TempDir(),
	})

	cfg := &Config{}
	err := cfg.loadETagFile()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
