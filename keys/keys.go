package keys

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

type List struct {
	Count int
	Data  []byte
}

func New(data []byte) *List {
	return &List{Data: data}
}

func (l *List) Validate() error {
	count := 0
	scanner := bufio.NewScanner(bytes.NewReader(l.Data))

	for number := 1; scanner.Scan(); number++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(line))
		if err != nil {
			return fmt.Errorf("parsing key on line %d: %w", number, err)
		}
		count++
	}

	err := scanner.Err()
	if err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}
	if count == 0 {
		return errors.New("reading keys: no valid key found")
	}
	l.Count = count

	return nil
}

func (l *List) Bytes() []byte {
	if bytes.HasSuffix(l.Data, []byte("\n")) {
		return l.Data
	}
	return append(l.Data, '\n')
}
