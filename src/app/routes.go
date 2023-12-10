package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/pedro-git-projects/chatbot-back/src/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *Application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/auth/signup", app.createUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/auth/signin", app.signinUserHandler)

	router.Handle(http.MethodGet, "/v1/user", app.jwtMiddleware(http.HandlerFunc(app.getUserHandler)))
	router.Handle(http.MethodPut, "/v1/user", app.jwtMiddleware(http.HandlerFunc(app.updateUserHandler)))
	router.Handle(http.MethodPatch, "/v1/user", app.jwtMiddleware(http.HandlerFunc(app.updateUserHandler)))
	router.Handle(http.MethodDelete, "/v1/user", app.jwtMiddleware(http.HandlerFunc(app.deleteUserHandler)))

	router.HandlerFunc(http.MethodGet, "/v1/swag", app.serveSwagger)

	router.Handler(http.MethodGet, "/v1/swagger/*any", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:4000/v1/swag"),
	))
	return router
}
