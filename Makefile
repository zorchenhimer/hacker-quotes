
#SOURCES= *.go \
#		 api/*.go \
#		 business/*.go \
#		 cmd/*.go \
#		 database/*.go \
#		 frontend/*.go \
#		 models/*.go

SOURCES=$(shell find . -type f -name "*.go")

bin/server: bin/ $(SOURCES)
	go build -o bin/server cmd/server.go

bin/:
	mkdir -p bin
