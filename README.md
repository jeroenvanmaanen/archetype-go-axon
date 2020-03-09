# Archetype of a Go project that uses Axon for Event Sourcing and CQRS

This project is still **Work In Progress**. It was cloned from [archetype-nix-go](https://github.com/jeroenvanmaanen/archetype-nix-go).
First point to solve:

The next phase is to communicate with Axon Server properly:
1. ☑ Set up a session
   *  ☑ Enable React app to call a RPC endpoint on the example-command-api service through grpc-web
2. ☑ Issue commands
3. ☐ Register a command handler and handle commands
4. ☐ Submit events
   * ☐ Stream events to UI
5. ☐ Retrieve the events for an aggregate and build a projection
   * ☐ Validate commands against the projection
6. ☐ Register a tracking event processor and handle events
7. ☐ Store records in a query model (Elastic Search?)
8. ☐ Register a query handler and handle queries
   * ☐ Show query results in UI

After that:

* Cache projections
* Store snapshots
* Use TLS
* Add claim-based security

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

To work with this project, you need to install docker.

Then, open a terminal (on windows either use git-bash or expose the  
docker daemon inside a docker container that has bash) and run:
```
[host]$ src/bin/clobber-build-and-run.sh --dev
```
The `--dev` flag specifies a configuration that runs React in development-mode
rather than optimized production mode. This makes fiddling with the UI much more straight-forward.

Then point your browser at
```
http://localhost:8123
```
Specify:
```
Grpc-swagger Server: localhost:8123
Endpoint Register: example-command-api:8181
```
Click "Register". Then a link appears under services:
org.leialearns.grpc.example.GreeterService. Click it and try the
method `/org.leialearns.grpc.example.GreeterService.Greet`.

## Setup step-by-step

Again, the main prerequisite is docker and either an extracted ZIP or a
clone of this project.

The first step after
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

## Run

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

## To do

Separate Docker image with build tools from Docker image for running the
compiled commands.

Add steps to generate te Go and JS stubs from the protocol buffer specifications to the Nix build.

I generated Go stubs for axon-server as follows:
```
[host]$ docker run --rm -ti -v "${HOME}:${HOME}" -w "$(pwd)" jeroenvm/archetype-nix-go bash
[container]# go get -u github.com/golang/protobuf/protoc-gen-go
[container]# go get google.golang.org/grpc
[container]# PATH="$PATH:/root/go/bin"
[container]# cd /src/axon-server-api/src/main/proto
[container]# bash WORKING_AREA/archetype-nix-go/src/bin/generate-proto-package.sh
```

Likewise I generated JS stubs for the example service:
```
[host]$ docker run --rm -ti -v "${HOME}:${HOME}" -w "$(pwd)" jeroenvm/archetype-nix-go bash
[container]# mkdir -p ~/home/bin
[container]# PATH="${PATH}:${HOME}/bin"
[container]# curl -L -sS -D - https://github.com/grpc/grpc-web/releases/download/1.0.7/protoc-gen-grpc-web-1.0.7-linux-x86_64 -o ~/bin/protoc-gen-grpc-web
[container]# chmod a+x ~/bin/protoc-gen-grpc-web
[container]# bash WORKING_AREA/archetype-nix-go/src/bin/generate-proto-js-package.sh
```