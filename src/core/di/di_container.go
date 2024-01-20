package di

import (
	. "enchantech-codex/src/feeds/controller"
	. "enchantech-codex/src/feeds/repository"
	. "enchantech-codex/src/feeds/service"
	. "enchantech-codex/src/users/controller"
	. "enchantech-codex/src/users/repository"
	. "enchantech-codex/src/users/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Container struct {
	EchoInstance *echo.Echo

	FeedRepository *FeedRepository
	FeedService    *FeedService
	FeedController *FeedController

	UserRepository *UserRepository
	UserService    *UserService
	UserController *UserController
}

func NewContainer(db *gorm.DB) *Container {
	echoInstance := echo.New()

	feedRepository := NewFeedRepository(db)
	feedService := NewFeedService(feedRepository)
	feedController := NewFeedController(feedService)

	userRepository := NewUserRepository(db)
	userService := NewUserService(userRepository)
	userController := NewUserController(userService)

	return &Container{
		EchoInstance:   echoInstance,
		FeedRepository: feedRepository,
		FeedService:    feedService,
		FeedController: feedController,
		UserRepository: userRepository,
		UserService:    userService,
		UserController: userController,
	}
}

func (c *Container) GetEchoInstance() *echo.Echo {
	return c.EchoInstance
}

func (c *Container) GetFeedRepository() *FeedRepository {
	return c.FeedRepository
}

func (c *Container) GetFeedService() *FeedService {
	return c.FeedService
}

func (c *Container) GetFeedController() *FeedController {
	return c.FeedController
}

func (c *Container) GetUserRepository() *UserRepository {
	return c.UserRepository
}

func (c *Container) GetUserService() *UserService {
	return c.UserService
}

func (c *Container) GetUserController() *UserController {
	return c.UserController
}
