// Filename: cmd/api/routes.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comments", a.requireActivatedUser(a.listCommentsHandler))

	router.HandlerFunc(http.MethodPost, "/v1/comments", a.requireActivatedUser(a.createCommentHandler))

	router.HandlerFunc(http.MethodGet, "/v1/comments/:id", a.requireActivatedUser(a.displayCommentHandler))

	router.HandlerFunc(http.MethodPatch, "/v1/comments/:id", a.requireActivatedUser(a.updateCommentHandler))

	router.HandlerFunc(http.MethodDelete, "/v1/comments/:id", a.requireActivatedUser(a.deleteCommentHandler))

	router.HandlerFunc(http.MethodPut, "/v1/users/activated", a.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", a.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", a.registerUserHandler)

	return a.recoverPanic(a.rateLimit(a.authenticate(router)))

}
