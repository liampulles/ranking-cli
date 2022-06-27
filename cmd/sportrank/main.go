package main

import (
	"os"

	"github.com/liampulles/span-digital-ranking-cli/cmd/sportrank/internal/wire"
)

func main() {
	os.Exit(wire.Run(os.Args, os.Stdin, os.Stdout))
}
