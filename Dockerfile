from golang:alpine as builder

run apk update && apk add --no-cache git
workdir /app
copy . .
run go mod tidy
run go build -o /hacker-quotes /app/cmd/server.go

from scratch
copy --from=builder /hacker-quotes /hacker-quotes
entrypoint ["/hacker-quotes"]
