package main

import (
    "fmt"
    submathpackage "github.com/jeroenvm/archetype-nix-go/src/pkg/submathpackage"
)

func main() {
	fmt.Println(submathpackage.Add(1,2))
}