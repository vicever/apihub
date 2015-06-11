package api_new

import (
	"net/http"

	"github.com/backstage/backstage/errors"
	"github.com/gorilla/context"
)

func (api *Api) requestIdMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
}

func (api *Api) authorizationMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authorization := r.Header.Get("Authorization")

	user, err := api.auth.UserFromToken(authorization)
	if err != nil {
		AddRequestError(r, errors.NewUnauthorizedError(errors.ErrLoginRequired))
		return
	}

	SetCurrentUser(r, user)
	next(rw, r)
}

func (api *Api) errorMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)

	err, ok := GetRequestError(r)
	if ok {
		handleError(rw, err)
		return
	}
}

func (api *Api) notFoundHandler(rw http.ResponseWriter, r *http.Request) {
	err := errors.NewNotFoundErrorNEW(errors.ErrNotFound)
	handleError(rw, &err)
}

func (api *Api) contextClearerMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer context.Clear(r)
	next(w, r)
}