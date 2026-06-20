APP=minilink

.PHONY: build test serve export icons

build:
	go build -trimpath -ldflags="-s -w" -o bin/$(APP) ./cmd/minilink

test:
	go test ./...

serve:
	go run ./cmd/minilink serve -config examples/mini-link.yaml

export:
	go run ./cmd/minilink export -config examples/mini-link.yaml -out dist

icons:
	go run ./cmd/minilink icons -out docs/icons.html
