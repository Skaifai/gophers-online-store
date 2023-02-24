package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.authMiddleware(app.healthcheckHandler))

	router.HandlerFunc(http.MethodPost, "/v1/auth/register", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/auth/authenticate", app.authenticateUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/auth/activate/:uuid", app.activateUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/auth/logout", app.authMiddleware(app.logoutUserHandler))
	router.HandlerFunc(http.MethodGet, "/v1/auth/refresh", app.authMiddleware(app.refreshHandler))

	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.showUserHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", app.updateUserHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.deleteUserHandler)

	return router
}
