// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"kratos-demo/internal/dao"
	"kratos-demo/internal/service"
	"kratos-demo/internal/server/grpc"
	"kratos-demo/internal/server/http"

	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.Provider, http.New, grpc.New, NewApp))
}
