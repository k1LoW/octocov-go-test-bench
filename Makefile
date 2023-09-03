PKG = github.com/k1LoW/octocov-go-test-bench
COMMIT = $$(git describe --tags --always)
OSNAME=${shell uname -s}
ifeq ($(OSNAME),Darwin)
	DATE = $$(gdate --utc '+%Y-%m-%d_%H:%M:%S')
else
	DATE = $$(date --utc '+%Y-%m-%d_%H:%M:%S')
endif

export GO111MODULE=on

BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)

default: test

ci: depsdev test test-integration

test:
	go test ./... -coverprofile=coverage.out -covermode=count

test-integration:
	go test -bench Benchmark -benchtime=1ms -benchmem -mod=mod ./testdata/bench -count 2 | go run ./cmd/octocov-go-test-bench/main.go > custom_metrics_bench.json

lint:
	golangci-lint run ./...

build:
	go build -ldflags="$(BUILD_LDFLAGS)" -o octocov-go-test-bench cmd/octocov-go-test-bench/main.go

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@latest
	go install github.com/Songmu/gocredits/cmd/gocredits@latest

prerelease_for_tagpr:
	gocredits -skip-missing -w
	cat _EXTRA_CREDITS >> CREDITS
	git add CHANGELOG.md CREDITS go.mod go.sum

release:
	git push origin main --tag
	goreleaser --clean

.PHONY: default test
