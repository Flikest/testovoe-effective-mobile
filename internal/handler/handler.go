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

	v1 := router.Group("/users")
	{
		v1.GET("")
		v1.DELETE("")
		v1.PATCH("")
		v1.POST("")
	}

	return router
}
