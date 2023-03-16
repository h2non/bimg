.PHONY: test clean examples

test:
	CGO_CFLAGS_ALLOW=-Xpreprocessor go test -cover .

clean:
	rm -rf ./bin

build:
	CGO_CFLAGS_ALLOW=-Xpreprocessor go build .

examples:
	CGO_CFLAGS_ALLOW=-Xpreprocessor go build -o bin/examples ./examples
	./bin/examples
