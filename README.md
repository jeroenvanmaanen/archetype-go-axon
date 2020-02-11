# Archetype of a Go project that can be built with Nix

This project is still **Work In Progress**. I think I've got most of
the necessary dependencies lined up, but I still have to connect the  
dots...

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

After that, build the executables from the Go code:
```
[host]$ src/bin/nix-build.sh
```

This is a list of commands that I figured out and might come in handy:
```
[host]$ docker run --rm -ti -v "${HOME}:${HOME}" -w "$(pwd)" jeroenvm/archetype-nix-go
[container]$ nix-shell -p "haskellPackages.ghcWithPackages (pkgs: [pkgs.http2-grpc-proto-lens])"
[nix-shell]# ghci
Prelude> import Network.GRPC.HTTP2.ProtoLens

[container]$ nix-shell --pure shell.nix
[nix-shell]# ghci
Prelude> import Network.HTTP2.Client

[container]$ nix-env -f '<.>' -iA haskellPackages."http2-grpc-proto-lens"
[container]$ nix-env -f '<.>' -iA haskellPackages.hoogle
[container]$ hoogle generate --insecure

[container]$ nix-shell --pure shell.nix --run "cabal repl"
```
