all: build

build: 
	go build -o lotus ./cmd/lotus

.PHONY: all build
