.PHONY: all clean
all: bin/obscure bin/obscure-cuda

bin/obscure: $(shell find cmd internal -name '*.go') go.mod
	mkdir -p bin
	go build -o $@ ./cmd/obscure

clean:
	rm -f bin/obscure bin/obscure-cuda
