BINARY  := resgate
MAIN    := ./main.go
VERSION := 1.0.0

.PHONY: build test install clean lint tidy

build:
	go build -o $(BINARY) $(MAIN)

test:
	go test ./...

install:
	go install ./...

clean:
	rm -f $(BINARY) $(BINARY).exe

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

demo: build
	./$(BINARY) pool create --name pool1 --cpu 32 --memory 65536 --gpu 8
	./$(BINARY) tenant add --name tenant1 --priority 5
	./$(BINARY) tenant add --name tenant2 --priority 8
	./$(BINARY) reserve --tenant tenant1 --pool pool1 --cpu 4 --memory 8192
	./$(BINARY) list
	./$(BINARY) status

.DEFAULT_GOAL := build
