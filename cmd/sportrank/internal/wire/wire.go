package wire

import "github.com/liampulles/span-digital-ranking-cli/cmd/sportrank/internal/driver/cli"

func Wire() cli.Engine {
	return cli.NewEngineImpl()
}
