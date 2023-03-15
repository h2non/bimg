.PHONY: test clean

test:
	CGO_CFLAGS_ALLOW=-Xpreprocessor go test -cover .

clean:
	rm -rf ./bin

build:
	CGO_CFLAGS_ALLOW=-Xpreprocessor go build .

resize:
	CGO_CFLAGS_ALLOW=-Xpreprocessor go build -o bin/resize ./examples/resize
	./bin/resize
