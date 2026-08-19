package config

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

type userLookup interface {
	Current() (*user.User, error)
}

type systemUser struct{}

func (systemUser) Current() (*user.User, error) {
	return user.Current()
}

type Config struct {
	File      string
	StateFile string
	URL       string

	user userLookup
}

func Load() (*Config, error) {
	cfg := &Config{user: systemUser{}}

	err := cfg.load()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) load() error {
	err := c.loadURL()
	if err != nil {
		return err
	}

	err = c.loadFile()
	if err != nil {
		return err
	}

	return c.loadStateFile()
}

func (c *Config) loadURL() error {
	url := os.Getenv("SSH_KEYS_URL")
	if url != "" {
		c.URL = url
		return nil
	}

	username, err := c.username()
	if err != nil {
		return err
	}
	c.URL = fmt.Sprintf("https://github.com/%s.keys", username)

	return nil
}

func (c *Config) username() (string, error) {
	username := os.Getenv("SSH_KEYS_USER")
	if username != "" {
		return username, nil
	}

	current, err := c.user.Current()
	if err != nil {
		return "", fmt.Errorf("loading config: %w", err)
	}
	if current.Username == "" {
		return "", errors.New("loading config: current user has no username")
	}

	return current.Username, nil
}

func (c *Config) loadFile() error {
	file := os.Getenv("SSH_KEYS_FILE")
	if file != "" {
		c.File = file
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	c.File = filepath.Join(home, ".ssh", "authorized_keys")

	return nil
}

func (c *Config) loadStateFile() error {
	if c.File == "" {
		return errors.New("loading config: keys file is not set")
	}

	c.StateFile = c.File + ".ssh-keys"

	return nil
}
