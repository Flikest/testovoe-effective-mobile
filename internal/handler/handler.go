package handler

import (
	"github.com/Flikest/testovoe-effective-mobile/internal/service"
	"github.com/gin-gonic/gin"
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

	return router
}
