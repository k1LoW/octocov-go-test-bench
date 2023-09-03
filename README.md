# octocov-go-test-bench [![CI](https://github.com/k1LoW/octocov-go-test-bench/actions/workflows/ci.yml/badge.svg)](https://github.com/k1LoW/octocov-go-test-bench/actions/workflows/ci.yml) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/octocov-go-test-bench/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/octocov-go-test-bench/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/octocov-go-test-bench/time.svg)

Generate [custom metrics JSON](https://github.com/k1LoW/octocov#custom-metrics) from the output of `go test -bench`.

## Usage

```console
$ go test -bench . -benchmem | octocov-go-test-bench
```

## Install

**go install:**

```console
$ go install github.com/k1LoW/octocov-go-test-bench/cmd/octocov-go-test-bench@latest
```

**deb:**

``` console
$ export OCTOCOV_GO_TEST_BENCH_VERSION=X.X.X
$ curl -o octocov-go-test-bench.deb -L https://github.com/k1LoW/octocov-go-test-bench/releases/download/v$OCTOCOV_GO_TEST_BENCH_VERSION/octocov-go-test-bench_$OCTOCOV_GO_TEST_BENCH_VERSION-1_amd64.deb
$ dpkg -i octocov-go-test-bench.deb
```

**RPM:**

``` console
$ export OCTOCOV_GO_TEST_BENCH_VERSION=X.X.X
$ yum install https://github.com/k1LoW/octocov-go-test-bench/releases/download/v$OCTOCOV_GO_TEST_BENCH_VERSION/octocov-go-test-bench_$OCTOCOV_GO_TEST_BENCH_VERSION-1_amd64.rpm
```

**apk:**

``` console
$ export OCTOCOV_GO_TEST_BENCH_VERSION=X.X.X
$ curl -o octocov-go-test-bench.apk -L https://github.com/k1LoW/octocov-go-test-bench/releases/download/v$OCTOCOV_GO_TEST_BENCH_VERSION/octocov-go-test-bench_$OCTOCOV_GO_TEST_BENCH_VERSION-1_amd64.apk
$ apk add octocov-go-test-bench.apk
```

**homebrew tap:**

```console
$ brew install k1LoW/tap/octocov-go-test-bench
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/octocov-go-test-bench/releases)
