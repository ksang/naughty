# Define parameters
BINARY=naughty
SHELL := /bin/bash
GOPACKAGES = $(shell go list ./... | grep ksang)
ROOTDIR = $(pwd)

.PHONY: build install test linux get-deps

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
	
get-deps:
	go get github.com/fatih/color
	go get github.com/mattn/go-colorable
	go get github.com/mattn/go-isatty

clone-deps:
	cd ${ROOTDIR}/vendor/src/github.com/fatih
	git clone https://github.com/fatih/color.git
	cd ${ROOTDIR}/vendor/src/github.com/mattn
	git clone https://github.com/mattn/go-colorable.git
	git clone https://github.com/mattn/go-isatty.git