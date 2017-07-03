# Server side application

## Requirements

[Installed Golang](https://golang.org/dl/)
[Runing PostgreSQL server](https://www.postgresql.org/docs/9.1/static/server-start.html)

## Building

Run `go build` in server dir.
After that command `server` executable must be available in bg/server dir.

## Running

Client side app should be built for this moment, so the server can serve static files (js, css and etc.)
(Read /bd/ui/README.md)

`cd bg` (project root directory)
Run `./server/server --help` to get available args.
