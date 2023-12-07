package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/auth/signup", app.createUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/auth/signin", app.signinUserHandler)
	router.Handle(http.MethodGet, "/v1/user", app.jwtMiddleware(http.HandlerFunc(app.getUserHandler)))

	return router
}
