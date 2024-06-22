# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./src
BINARY_NAME := query-benchmark

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	git diff --exit-code


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build
build:
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the  application (optional use `run csv=filepath` to specify a input file).
.PHONY: run
run: build
 ifdef csv
	/tmp/bin/${BINARY_NAME} -csv ${csv}
else
	/tmp/bin/${BINARY_NAME} 
endif

# ==================================================================================== #
# INFRA
# ==================================================================================== #
## start-timescaledb: run TimescaleDB in a container
.PHONY: start-timescaledb
start-timescaledb:
	cd timescaledb/ && make start-timescaledb

## stop-timescaledb: stop TimescaleDB in container
.PHONY: stop-timescaledb
stop-timescaledb:
	cd timescaledb/ && make stop-timescaledb