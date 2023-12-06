package main

import (
	"errors"
	"net/http"

	"github.com/pedro-git-projects/chatbot-back/internal/data/users"
	"github.com/pedro-git-projects/chatbot-back/internal/validator"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := users.CreateUserDTO{}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &users.User{
		Email:    payload.Email,
		Password: payload.Password,
		Name:     payload.Name,
		ImageURL: payload.ImageURL,
		Role:     payload.Role,
	}

	v := validator.New()
	payload.Validate(v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			app.badRequestResponse(w, r, errors.New("Credenciais tomadas"))
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, map[string]any{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

}
