//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/OICjangirrahul/students/internal/adapters/http"
	"github.com/OICjangirrahul/students/internal/adapters/repositories"
	"github.com/OICjangirrahul/students/internal/config"
	"github.com/OICjangirrahul/students/internal/core/ports"
	"github.com/OICjangirrahul/students/internal/core/services"
	"github.com/google/wire"
)

var repositorySet = wire.NewSet(
	repositories.NewSQLiteRepository,
	wire.Bind(new(ports.StudentRepository), new(*repositories.SQLiteRepository)),
)

func InitializeStudentHandler(cfg *config.Config) (*http.StudentHandler, error) {
	wire.Build(
		repositorySet,
		services.NewStudentService,
		http.NewStudentHandler,
	)
	return nil, nil
}
