# Define parameters
BINARY=naughty
SHELL := /bin/bash
GOPACKAGES = $(shell go list ./... | grep ksang)

.PHONY: build install test linux

GOPATH := ${PWD}/vendor:${GOPATH}

export GOPATH

default: build

build: main.go 
	go build -v -o ./build/${BINARY} main.go

install:
	go install  ./...

test:
	go test -race -cover ${GOPACKAGES}

clean:
	rm -rf build

linux: main.go
	GOOS=linux GOARCH=amd64 go build -o ./build/linux/${BINARY} main.go
