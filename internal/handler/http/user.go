package http

import (
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
