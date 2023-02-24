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
	router.HandlerFunc(http.MethodGet, "/v1/auth/activate/:uuid", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/auth/logout", app.logoutUserHandler)

	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.showUserHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", app.updateUserHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.deleteUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/products", app.addProductHandler)
	router.HandlerFunc(http.MethodGet, "/v1/products", app.listProductsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/products/:id", app.showProductHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/products/:id", app.deleteProductHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/products/:id", app.updateProductHandler)

	router.HandlerFunc(http.MethodPost, "/v1/cart/:id", app.listItemsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/cart-items", app.addItemToSessionHandler)
	router.HandlerFunc(http.MethodGet, "/v1/cart-items/:id", app.showItemHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/cart-items/:id", app.updateItemInSessionHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/cart-items/:id", app.removeItemFromSessionHandler)

	return router
}
