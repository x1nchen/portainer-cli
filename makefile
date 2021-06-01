GO = go
BUILD := ${GO} build
BIN := portainer-cli

PACKAGE_TEST := $(shell go list ./...|grep -v /test/integration/ | grep -v /test/e2e/)
PACKAGE_LIST  := go list ./...
PACKAGE_DIRECTORIES := $(PACKAGE_LIST)
FILES     := $$(find $$($(PACKAGE_DIRECTORIES)) -name "*.go")

REPO = $(shell if [ $${GITHUB_REPOSITORY+x} ]; then echo $$GITHUB_REPOSITORY; else basename `git rev-parse --show-toplevel`; fi)
BRANCH = $(shell if [ $${CI_HEAD_REF+x} ]; then echo $$CI_HEAD_REF; else basename `git symbolic-ref HEAD`; fi)
COMMIT = $(shell if [ $${CI_SHA_SHORT+x} ]; then echo $$CI_SHA_SHORT; else git rev-parse --short HEAD; fi)
VERSION = $(shell if [ $${CI_REF_NAME+x} ]; then echo $$CI_REF_NAME; else basename `git tag --sort=-version:refname | head -n 1`; fi)
VERSION_NAMESPACE := github.com/x1nchen/portainer-cli/cmd
BUILD_VERSION := ${BUILD} -ldflags "-X $(VERSION_NAMESPACE).Repo=$(REPO) -X $(VERSION_NAMESPACE).Branch=$(BRANCH) -X $(VERSION_NAMESPACE).Commit=$(COMMIT) -X $(VERSION_NAMESPACE).Version=$(VERSION)"
GOBUILD := ${BUILD_VERSION} -o ./bin/${BIN} ./main.go

all: build

dep:
	go mod download

build:
	${GOBUILD}

.PHONY: fmt
fmt:
	@echo "goimports with all source files"
	@goimports -w .

.PHONY: lint
lint:
	@echo skip go vet
	@golangci-lint run --deadline=5m

.PHONY: cov
cov: check
	gocov test $(packages) | gocov-html > coverage.html

.PHONY: check
check:
	@echo skip go vet

.PHONY: test-unit
test-unit:
	go test -short -v -coverprofile .coverage.txt $(shell go list ./... | grep -v /test/integration/)
	go tool cover -func .coverage.txt
