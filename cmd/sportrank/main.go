package main

import (
	"os"

	"github.com/liampulles/ranking-cli/cmd/sportrank/internal/wire"
)

func main() {
	engine := wire.Wire()
	os.Exit(engine.Run(os.Args, os.Stdin, os.Stdout))
}
