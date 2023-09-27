package user

import (
	"github.com/protomem/gotube/internal/database/postgres"
	handlhttp "github.com/protomem/gotube/internal/module/user/handler/http"
	"github.com/protomem/gotube/internal/module/user/repository"
	repopostgres "github.com/protomem/gotube/internal/module/user/repository/postgres"
	"github.com/protomem/gotube/internal/module/user/service"
	"github.com/protomem/gotube/internal/passhash/bcrypt"
	"github.com/protomem/gotube/pkg/logging"
)

type Module struct {
	*handlhttp.UserHandler
	service.UserService
	repository.UserRepository
}

func New(logger logging.Logger, db *postgres.DB) *Module {
	logger = logger.With("module", "user")

	userRepo := repopostgres.NewUserRepository(logger, db)
	userServ := service.NewUserService(bcrypt.NewHasher(bcrypt.DefaultCost), userRepo)
	userHandl := handlhttp.NewUserHandler(logger, userServ)

	return &Module{
		UserHandler:    userHandl,
		UserService:    userServ,
		UserRepository: userRepo,
	}
}
