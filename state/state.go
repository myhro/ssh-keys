package state

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/myhro/ssh-keys/store"
)

type Data struct {
	ETag   string `json:"etag"`
	SHA256 string `json:"sha256"`
}

type File struct {
	Path string

	store *store.File
}

func New(path string) *File {
	return &File{
		Path:  path,
		store: store.New(path),
	}
}

func (f *File) Read() (*Data, error) {
	raw, err := f.store.Read()
	if err != nil {
		return nil, err
	}

	data := &Data{}
	if len(raw) == 0 {
		return data, nil
	}

	err = json.Unmarshal(raw, data)
	if err != nil {
		slog.Warn("ignoring unreadable state file", "file", f.Path, "error", err)
		return &Data{}, nil
	}

	return data, nil
}

func (f *File) Write(data *Data) error {
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding %s: %w", f.Path, err)
	}

	return f.store.Write(append(raw, '\n'))
}

func Digest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
