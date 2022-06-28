package wire

import (
	"github.com/liampulles/span-digital-ranking-cli/cmd/sportrank/internal/adapter"
	"github.com/liampulles/span-digital-ranking-cli/cmd/sportrank/internal/driver/cli"
	"github.com/liampulles/span-digital-ranking-cli/cmd/sportrank/internal/usecase"
)

// TODO: Try wiring tool

func Wire() *cli.EngineImpl {
	usecaseSvc := usecase.NewServiceImpl()

	rowIOGateway := adapter.NewRowIOGatewayImpl(usecaseSvc)

	return cli.NewEngineImpl(rowIOGateway)
}
