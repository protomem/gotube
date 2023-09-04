package user

import (
	"github.com/protomem/gotube/internal/database"
	handlhttp "github.com/protomem/gotube/internal/module/user/handler/http"
	"github.com/protomem/gotube/internal/module/user/repository"
	"github.com/protomem/gotube/internal/module/user/repository/postgres"
	"github.com/protomem/gotube/internal/module/user/service"
	"github.com/protomem/gotube/internal/passhash/bcrypt"
	"github.com/protomem/gotube/pkg/logging"
)

type Module struct {
	*handlhttp.UserHandler
	service.UserService
	repository.UserRepository
}

func New(logger logging.Logger, db *database.DB) *Module {
	logger = logger.With("module", "user")

	userRepo := postgres.NewUserRepository(logger, db)
	userServ := service.NewUserService(bcrypt.NewHasher(bcrypt.DefaultCost), userRepo)
	userHandl := handlhttp.NewUserHandler(logger, userServ)

	return &Module{
		UserHandler:    userHandl,
		UserService:    userServ,
		UserRepository: userRepo,
	}
}
