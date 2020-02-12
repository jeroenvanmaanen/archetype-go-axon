package main

import (
    "fmt"
    submathpackage "github.com/jeroenvm/archetype-nix-go/src/pkg/submathpackage"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
)

var _ axonserver.PlatformOutboundInstruction

func main() {
	fmt.Println(submathpackage.Add(1,2))
}