# Hacker Quotes Generator

A generator for completely incorrect hacker quotes.  This is a reimplementation
of [GlOwl/hackerman](https://github.com/GlOwl/hackerman) in Go instead of
python.

## Build requirements

- Go 1.16 or newer
- GNU Make (optional)

To build just type `make` in the base directory and both the web server and
standalone command will be built.  If GNU Make isn't installed, running
`go build cmd/server.go` or `go build cmd/generate.go` will also work.

## Server Settings

The defualt settings for the web server are located in
`cmd/settings_default.json`.  Currently, the only database type supported is
sqlite.  To use an in-memory database append `?mode=memory` to the connection
string.

# TODO

- Make the webpage look less terrible.
- Write some documentation for the API.
- Write a twitter bot that uses the above API.

# License

License is MIT.  See `LICENSE.md` for details.
