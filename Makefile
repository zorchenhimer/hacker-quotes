
SOURCES=$(shell find . -not -path "./cmd/*" -type f -name "*.go")

all: bin/server bin/generate

bin/server: cmd/server.go bin/ $(SOURCES)
	go build -o bin/server $<

bin/generate: cmd/generate.go bin/ $(SOURCES)
	go build -o bin/server $<

bin/:
	mkdir -p bin

clean:
	rm -rf bin/
