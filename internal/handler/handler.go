package handler

import (
	_ "github.com/Flikest/testovoe-effective-mobile/api/docs"
	"github.com/Flikest/testovoe-effective-mobile/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{Service: s}
}

func (h Handler) InitRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		usersRouter := v1.Group("/user")
		usersRouter.GET("/", h.Service.GetUsers)
		usersRouter.DELETE("/:id", h.Service.DeleteUser)
		usersRouter.PATCH("/", h.Service.PatchUser)
		usersRouter.POST("/", h.Service.AddUser)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
