
SOURCES=$(shell find . -not -path "./cmd/*" -type f -name "*.go")
HTML=$(shell find frontend/ -type f -name "*.html")

all: bin/HackerServer bin/HackerQuote

bin/HackerServer: cmd/server.go $(SOURCES) $(HTML)
	go build -o $@ $<

bin/HackerQuote: cmd/generate.go $(SOURCES)
	go build -o $@ $<

clean:
	rm -rf bin/
