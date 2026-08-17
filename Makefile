BINARY = ssh-keys

.PHONY: build clean

build:
	go build -trimpath -ldflags="-s -w" -o $(BINARY) .

clean:
	rm -f $(BINARY)

version:
	go version -m ssh-keys
