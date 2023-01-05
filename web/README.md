# SimpleMixer - web

This directory contains all the front and back-end related files to run the main central server
and to run a simple front-end so users can interact with the contract directly.

## Deployment

Just the compilation of the Go application is needed using `go build`.
By default the HTTP server will listen on `:8080`.

First run:

```shell
go build
```

Then execute the compiled binary that is produced.
Visit `localhost:8080` to view the front-end.

The following flags should be provided when starting the application:

- privateKey: Hexadecimal string representation of the private key to use for ECDSA signatures.
- contractAddress: Address where the contract is currently deployed.

**The private key is passed raw as an argument for the PoC. Outside of the PoC it makes sense to use a more secure approach like an external signer**

## Testing

The following command will execute all the repository tests:

```shell
go test
```
