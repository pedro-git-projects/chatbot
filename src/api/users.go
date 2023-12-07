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

	token, err := app.generateJWT(user.ID, user.Role)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]any{
		"user":  user,
		"token": token,
	}

	err = app.writeJSON(w, http.StatusCreated, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) signinUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := users.LoginUserDTO{}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	payload.Validate(v)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.Authenticate(payload.Email, payload.Password)
	if err != nil {
		app.unauthorizedResponse(w, r, "Credenciais inválidas")
		return
	}

	token, err := app.generateJWT(user.ID, user.Role)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	if err := app.writeJSON(w, http.StatusOK, response, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		app.unauthorizedResponse(w, r, "ID do usuário não foi encontrado no contexto da requisição")
		return
	}

	user, err := app.models.Users.Get(userID)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, user, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		app.unauthorizedResponse(w, r, "ID do usuário não foi encontrado no contexto da requisição")
		return
	}

	payload := users.UpdateUserDTO{}

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

	updatedUser, err := app.models.Users.Update(userID, user)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			app.badRequestResponse(w, r, errors.New("Credenciais tomadas"))
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, updatedUser, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		app.unauthorizedResponse(w, r, "ID do usuário não foi encontrado no contexto da requisição")
		return
	}

	err := app.models.Users.Delete(userID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
