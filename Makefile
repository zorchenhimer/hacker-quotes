
SOURCES=$(shell find . -not -path "./cmd/*" -type f -name "*.go")
HTML=$(shell find frontend/ -type f -name "*.html")

all: bin/server bin/generate

bin/server: cmd/server.go $(SOURCES) $(HTML)
	go build -o $@ $<

bin/generate: cmd/generate.go $(SOURCES)
	go build -o $@ $<

clean:
	rm -rf bin/
