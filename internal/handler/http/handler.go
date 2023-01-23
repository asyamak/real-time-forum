package http

import (
	"real-time-forum/internal/service"

	"github.com/rshezarr/gorr"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *gorr.Router {
	router := gorr.NewRouter()

	router.POST("/auth/sign-up", h.SignUp)
	router.POST("/auth/sign-in", h.SignIn)

	return router
}
