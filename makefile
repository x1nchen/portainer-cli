GO = go
GOGET := $(GO) get
BUILD := ${GO} build ${GOFLAGS}
BIN := portainer-cli
GOBUILD := ${BUILD} -o ./bin/${BIN} ./main.go
PACKAGE_TEST := $(shell go list ./...|grep -v /test/integration/ | grep -v /test/e2e/)
PACKAGE_LIST  := go list ./...
PACKAGE_DIRECTORIES := $(PACKAGE_LIST)
FILES     := $$(find $$($(PACKAGE_DIRECTORIES)) -name "*.go")

all: build

.PHONY: lint
lint:
	@golangci-lint run --skip-dirs=test/integration --deadline=5m

.PHONY: dep
dep:
	go mod download

# gen mock or others
.PHONY: gen
gen:
	go generate ./...

# gen mock or others
.PHONY: gen-mock
gen-mock:
	mockgen -source=session/session.go -destination=session/session_mock.go -package=session Manager

.PHONY: fmt
fmt:
	@echo "goimports"
	@goimports -w .

.PHONY: build
build:
	mkdir -p bin
	${GOBUILD}


.PHONY: test-unit
test-unit:
	# go test -race -v $(PACKAGE_TEST) -coverprofile .coverage.txt
	go test -v $(PACKAGE_TEST) -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

.PHONY: test-integration
test-integration:
	go test -v ./test/integration/...

.PHONY: docker-build
docker-build:
	mkdir -p bin
	docker run --rm -v ${PWD}:/usr/src/build -v ${GOPATH}/pkg/mod:/go/pkg/mod -w /usr/src/build ${GO_BUILD_IMAGE} ${GOBUILD}

.PHONY: docker
docker:
	docker build --rm -t $(IMAGE_NAME) .
