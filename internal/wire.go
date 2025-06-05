//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/OICjangirrahul/students/internal/adapters/http"
	"github.com/OICjangirrahul/students/internal/adapters/repositories"
	"github.com/OICjangirrahul/students/internal/core/ports"
	"github.com/OICjangirrahul/students/internal/core/services"
	"github.com/google/wire"
)

var dbSet = wire.NewSet(repositories.NewDB)

var studentRepositorySet = wire.NewSet(
	repositories.NewStudentRepository,
	wire.Bind(new(ports.StudentRepository), new(*repositories.StudentRepository)),
)

var teacherRepositorySet = wire.NewSet(
	repositories.NewTeacherRepository,
	wire.Bind(new(ports.TeacherRepository), new(*repositories.TeacherRepository)),
)

var studentServiceSet = wire.NewSet(
	services.NewStudentService,
)

var teacherServiceSet = wire.NewSet(
	services.NewTeacherService,
)

type Handlers struct {
	Student *http.StudentHandler
	Teacher *http.TeacherHandler
}

func InitializeHandlers() (*Handlers, error) {
	wire.Build(
		dbSet,
		studentRepositorySet,
		teacherRepositorySet,
		studentServiceSet,
		teacherServiceSet,
		http.NewStudentHandler,
		http.NewTeacherHandler,
		wire.Struct(new(Handlers), "*"),
	)
	return nil, nil
}
