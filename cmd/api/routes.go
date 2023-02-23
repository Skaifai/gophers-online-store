package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/auth/register", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/auth/authenticate", app.authenticateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/auth/logout", app.logoutUserHandler)

	router.HandlerFunc(http.MethodGet, "/v1/users/activate/:uuid", app.activateUserHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.deleteUserHandler)

	return router
}
