GO_VERSION:=$(shell go version)

.PHONY: bench profile clean

all: install

test:
	go test -v -race ./...

bench:
	go test -count=5 -run=NONE -bench . -benchmem

clean:
	rm -rf pprof
	rm -rf bench
