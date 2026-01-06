# httpfromtcp

Simple HTTP server implementation built from scratch using TCP sockets in Go.

## Run

```bash
go run cmd/httpserver/main.go
```

Server starts on port `42069`.

## Endpoints

- `/` - 200 OK
- `/yourproblem` - 400 Bad Request
- `/myproblem` - 500 Internal Server Error
