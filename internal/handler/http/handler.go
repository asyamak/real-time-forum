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
	// user handlers
	router.POST("/api/user/sign-up", h.SignUp)
	router.POST("/api/user/sign-in", h.SignIn)
	router.GET("/api/user/:user_id", h.GetUser)
	router.GET("/api/user/:user_id/posts", h.GetUserPosts)
	router.GET("/api/user/:user_id/liked-posts", h.GetUserVotedPosts)

	//post handlers

	//categories handlers

	//comments handlers

	//chat handlers

	//images fileserver

	return router
}
