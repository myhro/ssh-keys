package keys

import (
	"crypto/ed25519"
	"crypto/rand"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

func publicKey(t *testing.T) string {
	t.Helper()

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generating key: %v", err)
	}
	key, err := ssh.NewPublicKey(pub)
	if err != nil {
		t.Fatalf("converting key: %v", err)
	}

	return strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key)))
}

func TestValidate(t *testing.T) {
	first := publicKey(t)
	second := publicKey(t)

	tests := []struct {
		name  string
		data  string
		count int
		fail  bool
	}{
		{
			name:  "single key",
			data:  first + "\n",
			count: 1,
		},
		{
			name:  "multiple keys",
			data:  first + "\n" + second + "\n",
			count: 2,
		},
		{
			name:  "key with comment",
			data:  first + " user@host\n",
			count: 1,
		},
		{
			name:  "blank lines and comments",
			data:  "# header\n\n" + first + "\n\n",
			count: 1,
		},
		{
			name:  "no trailing newline",
			data:  first,
			count: 1,
		},
		{
			name: "empty data",
			fail: true,
		},
		{
			name: "only comments",
			data: "# nothing to see here\n",
			fail: true,
		},
		{
			name: "invalid key",
			data: "ssh-ed25519 not-a-real-key\n",
			fail: true,
		},
		{
			name: "one invalid key among valid ones",
			data: first + "\nssh-rsa AAAA\n" + second + "\n",
			fail: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			list := New([]byte(tc.data))
			err := list.Validate()

			if tc.fail {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				if list.Count != 0 {
					t.Fatalf("expected no keys counted, got %d", list.Count)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if list.Count != tc.count {
				t.Fatalf("expected %d keys, got %d", tc.count, list.Count)
			}
		})
	}
}

func TestBytes(t *testing.T) {
	key := publicKey(t)

	tests := []struct {
		name string
		data string
		want string
	}{
		{
			name: "adds missing newline",
			data: key,
			want: key + "\n",
		},
		{
			name: "keeps existing newline",
			data: key + "\n",
			want: key + "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := string(New([]byte(tc.data)).Bytes())
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}
