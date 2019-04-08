# galaxy
[![GoDoc](https://godoc.org/github.com/henvic/galaxy?status.svg)](https://godoc.org/github.com/henvic/galaxy) [![Build Status](https://travis-ci.org/henvic/galaxy.svg?branch=master)](https://travis-ci.org/henvic/galaxy) [![Coverage Status](https://coveralls.io/repos/henvic/galaxy/badge.svg)](https://coveralls.io/r/henvic/galaxy) [![codebeat badge](https://codebeat.co/badges/fc6e1e79-504b-4e49-9607-9767abdcf6c2)](https://codebeat.co/projects/github-com-henvic-galaxy-master) [![Go Report Card](https://goreportcard.com/badge/github.com/henvic/galaxy)](https://goreportcard.com/report/github.com/henvic/galaxy)

galaxy offers a drone navigation service (DNS).

## Using
Querying a DNS location on sectorID 311:

```bash
curl -XPOST http://localhost:9000/v1/sectors/311/dns -H "Content-Type: application/json" -d '{"x": "33", "y": "42", "z": "13", "vel": "4.229"}'
```

## Commands
Use `make server` to run the application or `make apidocs` to run the application with [Swagger](https://swagger.io) API documentation server turned on.

### Docker
A Dockerfile is available in this directory as well.

You can build and run it with

```bash
docker build -t galaxy:latest
docker run -p 9000:9000 galaxy
```

An image is not currently available at the Docker Hub registry.

### Ports and configuration
The application runs at port 9000 and also exposes the Go profiler and debugger at local port 8081.

You might want to run `cmd/server --help` for a list of available flags.

The environment variable `DEBUG` sets the logging to debug mode.

## Contributing
You can get the latest source code with `go get -u github.com/henvic/galaxy`

The following commands are available and require no arguments:

* **make test**: run tests

In lieu of a formal style guide, take care to maintain the existing coding style. Add unit tests for any new or changed functionality. Integration tests should be written as well.

On API changes, update the Swagger documentation accordingly using [swag](https://github.com/swaggo/swag).
Use `make write-apidocs` to recreate the API documentation.

## Committing and pushing changes
The master branch of this repository on GitHub is protected:
* force-push is disabled
* tests MUST pass on Travis before merging changes to master
* branches MUST be up to date with master before merging

Keep your commits neat and [well documented](https://wiki.openstack.org/wiki/GitCommitMessages). Try to always rebase your changes before publishing them.

## Maintaining code quality
[goreportcard](https://goreportcard.com/report/github.com/henvic/galaxy) can be used online or locally to detect defects and static analysis results from tools with a great overview.

Using go test and go cover are essential to make sure your code is covered with unit tests.

Always run `make test` before submitting changes.
