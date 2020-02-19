# Archetype of a Go project that uses Axon for Event Sourcing and CQRS

This project is still **Work In Progress**. It was cloned from [archetype-nix-go](https://github.com/jeroenvanmaanen/archetype-nix-go).
The next phase is to communicate with Axon Server properly:
1. Set up a session
2. Issue commands
3. Register a command handler and handle commands
4. Submit events
5. Register a tracking event processor and handle events
6. Store records in a query model (Elastic Search?)
7. Register a query handler and handle queries

## Introduction

My aim is to create a project that
can be used as a template for future projects with the following
characteristics:
* Command / Query Responsibility Segregation (CQRS)
* Event Sourcing
* Strong typing
* High Performance
* High Availability
* Scalability

On top of that, I would love to be able to start a project as a monolithic
application and have it evolve into a collection of micro-services that is
integrated into a service mesh architecture.

The components that I want to combine are:
* The Go language (for high-performance, type safety and gRPC integration with Axon Server)
* Docker (to minimise the impact on/from the host system)
* Nix (to manage dependencies)
* Axon Server (for event storage, message routing, and scalability)
* Envoy (for service mesh architecture and high availability)

## Setup

I mostly followed [Golang Demo](https://github.com/MatrixAI/Golang-Demo)
by _Roger Qiu_. Any flaws are of course my own.

To work with this project, you need to install docker. The first step after
that is to acquire a docker image that has Nix and Go tools. It will be pulled from
docker hub automatically the first time you run `docker run` or
`docker-compose up`. You can also build it yourself with:
```
[host]$ src/bin/docker-build.sh # Optional. It is also available on Docker Hub
```

The docker compose script that we are going to run later, needs some local settings.
Create them from provided sample files for now:
```
[host]$ src/bin/create-local-settings.sh
```

You might need to update `deps.nix`:
```
[host]$ docker run --rm -ti -v "${HOME}:${HOME}" -w "$(pwd)" jeroenvm/archetype-nix-go bash
[container]# vgo2nix
```

After that, build the executables from the Go code:
```
[host]$ src/bin/nix-build.sh
```

Now it is time to start two docker containers:
```
[host]$ src/docker/docker-compose-up.sh
```

This is a good time to open a browser window for the [AxonDashboard](http://localhost:8024)
and keep it in view when running the example.

In another terminal window, open a bash prompt in the container that can run the example and run it:
```
[host2]$ docker exec -ti -w "$(pwd)" example_example_1 bash
[container]# result/bin/example
```
During the built-in sleep of 10 seconds, a box labeled GoClient should pop up
in the AxonDashboard. Id disappears again when the example application stops.

Building executables also works inside the docker container:
```
[container]# src/bin/nix-build.sh
```

I generated Go stubs for axon-server as follows:
```
[host]$ docker run --rm -ti -v "${HOME}:${HOME}" -w "$(pwd)" jeroenvm/archetype-nix-go bash
[container]# go get -u github.com/golang/protobuf/protoc-gen-go
[container]# go get google.golang.org/grpc
[container]# PATH="$PATH:/root/go/bin"
[container]# cd /src/axon-server-api/src/main/proto
[container]# bash WORKING_AREA/archetype-nix-go/src/bin/generate-proto-package.sh
```
