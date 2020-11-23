.POSIX:
.SUFFIXES:
.PHONY: build clean deps install

PREFIX = /usr/local
INSTALLPATH = $(PREFIX)/bin/dorts

all: build

build:
	go build -o dorts -v

deps:
	go get github.com/pelletier/go-toml

install:build
	install dorts $(INSTALLPATH)

clean:
	go clean
