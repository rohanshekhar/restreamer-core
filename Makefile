COMMIT := $(shell if [ -d .git ]; then git rev-parse HEAD; else echo "unknown"; fi)
SHORTCOMMIT := $(shell echo $(COMMIT) | head -c 7)
BRANCH := $(shell if [ -d .git ]; then git rev-parse --abbrev-ref HEAD; else echo "master"; fi)
BUILD := $(shell date -u "+%Y-%m-%dT%H:%M:%SZ")
OSARCH := $(shell if [ "${GOOS}" -a "${GOARCH}" ]; then echo "-${GOOS}-${GOARCH}"; else echo ""; fi)

all: build

## build: Build core (default)
build:
	go build -o core$(OSARCH)

## swagger: Update swagger API documentation (requires github.com/swaggo/swag)
swagger:
	swag init -g http/server.go

## gqlgen: Regenerate GraphQL server from schema
gqlgen:
	go run github.com/99designs/gqlgen generate --config http/graph/gqlgen.yml

## test: Run all tests
test:
	go test -coverprofile=/dev/null ./...

## vet: Analyze code for potential errors
vet:
	go vet ./...

## fmt: Format code
fmt:
	go fmt ./...

## update: Update dependencies
update:
	go get -u
	@-$(MAKE) tidy

## tidy: Tidy up go.mod
tidy:
	go mod tidy

## vendor: Update vendored packages
vendor:
	go mod vendor

## run: Build and run core
run: build
	./core

## lint: Static analysis with staticcheck
lint:
	staticcheck ./...

## import: Build import binary
import:
	cd app/import && go build -o ../../import -ldflags="-s -w"

## coverage: Generate code coverage analysis
coverage:
	go test -coverprofile test/cover.out ./...
	go tool cover -html=test/cover.out -o test/cover.html

## commit: Prepare code for commit (vet, fmt, test)
commit: vet fmt lint test build
	@echo "No errors found. Ready for a commit."

## release: Build a release binary of core
release:
	go build -o core -ldflags="-s -w -X github.com/datarhei/core/app.Commit=$(COMMIT) -X github.com/datarhei/core/app.Branch=$(BRANCH) -X github.com/datarhei/core/app.Build=$(BUILD)"

## docker: Build standard Docker image
docker:
	docker build -t core:$(SHORTCOMMIT) .

.PHONY: help build swagger test vet fmt vendor commit coverage lint release import

## help: Show all commands
help: Makefile
	@echo
	@echo " Choose a command:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
