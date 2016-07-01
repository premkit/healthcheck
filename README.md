Healthcheck [![CircleCI](https://circleci.com/gh/premkit/healthcheck.svg?style=svg)](https://circleci.com/gh/premkit/healthcheck) [![Coverage Status](https://coveralls.io/repos/github/premkit/healthcheck/badge.svg?branch=master)](https://coveralls.io/github/premkit/healthcheck?branch=master)
===========

## Setup
We use Docker for the official build environment.  Docker is not required to run the binary, but to simplify the development and build enviroment, there are 
Dockerfiles to use if you want to build the binary.

## Build the development container
```shell
$ docker build -t premkit/healthcheck:dev .
```

## Run the development environment

### Start the container, build the execuatable, and run the service.
```
$ make shell
$ make build run
```

## Run tests
Tests are automatically run in CircleCI after pushing.  Tests can be run manually with
```shell
$ make shell
$ make test
```
