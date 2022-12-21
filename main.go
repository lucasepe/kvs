package main

import (
	"github.com/lucasepe/kvs/cmd"
)

// Build information. Populated at build-time.
var (
	Version string
	Build   string
)

func main() {
	cmd.Run(Version, Build)
}
