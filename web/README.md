# SimpleMixer - web

This directory contains all the front and back-end related files to run the main central server
and to run a simple front-end so users can interact with the contract directly.

## Testing

The following command will execute all the repository tests:

```shell
go test
```

## Deployment

Just compilation of the Go application is needed using `go build`.
By default the HTTP server will listen on `:8080`.

First run:

`go build`

Then execute the compiled binary that is produced.
Visit `localhost:8080` to view the front-end.
