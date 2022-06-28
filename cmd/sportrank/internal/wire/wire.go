package wire

import "github.com/liampulles/span-digital-ranking-cli/cmd/sportrank/internal/driver/cli"

// TODO: Try wiring tool

func Wire() cli.Engine {
	return cli.NewEngineImpl()
}
