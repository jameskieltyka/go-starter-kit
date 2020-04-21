#! /bin/bash
GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_TIME := $(shell date -Ins)
BUILD_VERSION := $(shell cat ./VERSION)

.PHONY: build
build:
	go build -ldflags "-X github.com/jkieltyka/go-starter-kit/internal/version.GitCommit=$(GIT_COMMIT) \
	-X github.com/jkieltyka/go-starter-kit/internal/version.BuildTime=$(BUILD_TIME) \
	-X github.com/jkieltyka/go-starter-kit/internal/version.BuildVersion=$(BUILD_VERSION)" .