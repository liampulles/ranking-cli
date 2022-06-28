package wire

import (
	"github.com/liampulles/ranking-cli/cmd/sportrank/internal/adapter"
	"github.com/liampulles/ranking-cli/cmd/sportrank/internal/driver/cli"
	"github.com/liampulles/ranking-cli/cmd/sportrank/internal/usecase"
)

func Wire() *cli.EngineImpl {
	usecaseSvc := usecase.NewServiceImpl()

	rowIOGateway := adapter.NewRowIOGatewayImpl(usecaseSvc)

	return cli.NewEngineImpl(rowIOGateway)
}
