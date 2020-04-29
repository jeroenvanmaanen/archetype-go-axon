# Archetype of a Go project that uses Axon for Event Sourcing and CQRS

## Introduction

This aims to be a project that can be used as a template for future
projects with the following characteristics:
* Command / Query Responsibility Segregation (CQRS)
* Event Sourcing
* Strong typing
* High Performance
* High Availability
* Scalability

On top of that, wouldn't it be fun to be able to start a project as a monolithic application and have it evolve into a collection of micro-services that is integrated into a service mesh architecture?

The components that this project aims to combine are:
* The Go language (for high-performance, type safety and gRPC integration with Axon Server)
* Docker (to minimise the impact on/from the host system)
* Nix (to manage dependencies)
* Axon Server (for event storage, message routing, and scalability)
* Envoy (for service mesh architecture and high availability)

## Status

The first phase is already done.
The current project communicates with Axon Server properly:
1. ☑ Set up a session
   *  ☑ Enable React app to call a RPC endpoint on the example-command-api service through grpc-web
2. ☑ Issue commands
3. ☑ Register a command handler and handle commands
4. ☑ Submit events
   * ☑ Stream events to UI
5. ☑ Retrieve the events for an aggregate and build a projection
   * ☑ Validate commands against the projection
6. ☑ Register a tracking event processor and handle events
7. ☑ Store records in a query model: Elastic Search
   * ☑ Store tracking token in Elastic Search
8. ☑ Register a query handler and handle queries
   * ☑ Show query results in UI

Other features:

* Configuration properties in event store
* Public key management in event store
  * Public key of initially trusted key manager compiled into binary
  * Private key to sign JWT tokens and decode credentials must be uploaded on startup with a challenge signed by a trusted key manager
* Claim-based security based on JWT
  * Todo: supply JWT in headers/metadata rather than payload

After that:

* Cache projections
* Store snapshots
* Use TLS
* Extract reusable code into a proper library that can be included as a dependency (evolve into a proper framework on par with AxonFramework for Java)
* Add context management with proper canceling of operations
* Support distributable segmented tracking event processors
* Fix bug with disappearing connections when Go applications talk to
  AxonServer via Envoy
* Separate Docker image with build tools from Docker image for running the
compiled commands
* Provide CMD or PowerShell script to run Docker-in-Docker for Windows users

This project started as a clone of [archetype-nix-go](https://github.com/jeroenvanmaanen/archetype-nix-go).
I mostly followed [Golang Demo](https://github.com/MatrixAI/Golang-Demo)
by _Roger Qiu_. Any flaws are of course my own.

## Quick start

To work with this project, you need to install docker.

Then, open a terminal (on windows either use git-bash or expose the docker daemon inside a docker container that has bash) and run:
```
[host]$ src/bin/clobber-build-and-run.sh --dev
```
The `--dev` flag specifies a configuration that runs React in development-mode rather than optimized production mode.
This makes fiddling with the UI much more straight-forward.

Then point your browser at
```
http://localhost:3000
```
to interact with the web front-end built in React.

It is possible to talk to the gRPC API directly using swagger-grpc:
```
http://localhost:8123
```
Specify:
```
Grpc-swagger Server: localhost:8123
Endpoint Register: example-command-api:8181
```
Click "Register".
Then a link appears under services: org.leialearns.grpc.example.GreeterService.
Click it and try the method `/org.leialearns.grpc.example.GreeterService.Greet`.

## Step-by-step

Again, the main prerequisite is docker and either an extracted ZIP or a
clone of this project.

The first step after that is to acquire a docker image that has Nix and Go tools.
It will be pulled from docker hub automatically the first time you run `docker run` or `docker-compose up`.
You can also build it yourself with:
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
[host]$ docker run --rm -ti -v "${HOME}:${HOME}" -w "$(pwd)" jeroenvm/build-protoc bash
[container]# vgo2nix
```

After that, build the executables from the Go code:
```
[container]$ src/bin/nix-build.sh
```
or
```
[container]$ exit
[host]$ src/bin/nix-build.sh
```
If any dependencies have changed, then the output of `nix-build.sh` will specify what the new value for key `modSha256` in `default.nix` has to be.
After `default.nix` has been updated, `nix-build.sh` should be run again.

### Run

Now it is time to start the docker containers:
```
[host]$ src/docker/docker-compose-up.sh
```
or, for easier front-end development:
```
[host]$ src/docker/docker-compose-up.sh --dev
```
This starts a number of docker containers:
* present: either Nginx (optimized) or Node Express (development) that serves the presentation layer
* command-api: executables compiled from Go with the business logic
* axon-server: event store and message routing
* elastic-search: persistence of query models
* proxy: Envoy proxy to manage network traffic for the server components, both between each other and with the public internet
* grpc-swagger: Swagger UI for gRPC API

This is a good time to open a browser window for the [AxonDashboard](http://localhost:8024) and keep it in view when running the example.

A few boxes labeled GoClient should pop up in the AxonDashboard.
Note that these boxes correspond to logical components (API, Command Handler, Event Processor, Query Handler) not server processes or containers.

Building executables also works inside the docker container:
```
[container]# src/bin/nix-build.sh
```