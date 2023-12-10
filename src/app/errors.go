package app

import (
	"fmt"
	"net/http"
)

func (app Application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

func (app Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := app.writeJSON(w, status, map[string]any{"erro": message}, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	msg := "O servidor encontrou um problema e não foi capaz de processar a sua requisição"
	app.errorResponse(w, r, http.StatusInternalServerError, msg)

}
func (app Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "O recurso solicitado não foi encontrado"
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Este recurso não suporta o método %s", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, map[string]any{"erro": msg})
}

func (app Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app Application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app Application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, message string) {
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
