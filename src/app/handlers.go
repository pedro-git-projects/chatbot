package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/pedro-git-projects/chatbot-back/src/data/users"
	"github.com/pedro-git-projects/chatbot-back/src/data/validator"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Ambiente string `json:"ambiente"`
	Versao   string `json:"versão"`
}

// @Summary Health Check
// @Description Retorna o status da aplicação.
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse "{'status': 'disponível', 'ambiente': 'dev', 'versão': '1.0'}"
// @Router /health [get]
func (app Application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := HealthResponse{
		Status:   "disponível",
		Ambiente: app.config.env,
		Versao:   app.config.version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) createUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) signinUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) serveSwagger(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("src/docs/swagger.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Decode JSON from the file
	var data interface{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	json.NewEncoder(w).Encode(data)
}
