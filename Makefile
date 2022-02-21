
SOURCES=$(shell find . -not -path "./cmd/*" -type f -name "*.go")
HTML=$(shell find frontend/ -type f -name "*.html")
#USER=data-www

all: bin/HackerServer bin/HackerQuote

bin/HackerServer: cmd/server.go $(SOURCES) $(HTML)
	go$(GO_VERSION) build -o $@ $<

bin/HackerQuote: cmd/generate.go $(SOURCES)
	go$(GO_VERSION) build -o $@ $<

clean:
	rm -rf bin/

# TODO: make this proper
install: bin/HackerServer bin/HackerQuote
	mkdir -p /opt/HackerQuotes
	cp bin/HackerServer /opt/HackerQuotes
	cp bin/HackerQuote /opt/HackerQuotes
	#chown -R $(USER):$(USER) /opt/HackerQuotes
	#cp HackerServer.service /etc/systemd/system/
