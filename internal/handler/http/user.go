package http

import (
	"errors"
	"net/http"
	"real-time-forum/internal/service"

	"github.com/rshezarr/gorr"
)

type usersSignUpInput struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}

func (h *Handler) SignUp(c *gorr.Context) {
	var input usersSignUpInput

	if err := c.ReadBody(&input); err != nil {
		c.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.User.SignUp(c.Context(), service.UserSignUpInput{
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Age:       input.Age,
		Gender:    input.Gender,
		Email:     input.Email,
		Password:  input.Password,
	}); err != nil {
		c.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	c.WriteHeader(http.StatusCreated)
}

type usersSignInInput struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (h *Handler) SignIn(c *gorr.Context) {
	var input usersSignInInput

	if err := c.ReadBody(&input); err != nil {
		c.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.User.SignIn(c.Context(), service.UserSignInInput{
		UsernameOrEmail: input.UsernameOrEmail,
		Password:        input.Password,
	})
	if err != nil {
		c.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	resp := tokenResponse{
		Token: token,
	}

	c.WriteJSON(http.StatusOK, resp)
}

func (h *Handler) GetUser(c *gorr.Context) {
	userID, err := c.GetIntParam("user_id")
	if err != nil {
		c.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.User.GetByID(c.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExists) {
			c.WriteError(http.StatusNotFound, err.Error())
			return
		}
		c.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	c.WriteJSON(http.StatusOK, user)
}

func (h *Handler) GetUserPosts(c *gorr.Context) {
	userID, err := c.GetIntParam("user_id")
	if err != nil {
		c.WriteError(http.StatusBadRequest, err.Error())
		return
	}

	posts, err := h.service.User.GetUsersPosts(c.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExists) {
			c.WriteError(http.StatusNotFound, err.Error())
			return
		}
		c.WriteError(http.StatusInternalServerError, err.Error())
		return
	}

	c.WriteJSON(http.StatusOK, posts)
}
