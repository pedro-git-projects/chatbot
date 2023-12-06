package main

import (
	"fmt"
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

	v := validator.New()
	payload.Validate(v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", payload)

}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

}
