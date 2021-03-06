package api

import (
	"encoding/json"
	"net/http"

	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/errors"
	"github.com/gorilla/mux"
)

func (api *Api) appCreate(rw http.ResponseWriter, r *http.Request, user *account.User) {
	app := account.App{}
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}

	team, err := findTeamAndCheckUser(app.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err := app.Create(*user, *team); err != nil {
		handleError(rw, err)
		return
	}

	Created(rw, app)
}

func (api *Api) appUpdate(rw http.ResponseWriter, r *http.Request, user *account.User) {
	app, err := account.FindAppByClientId(mux.Vars(r)["client_id"])
	if err != nil {
		handleError(rw, err)
		return
	}

	_, err = findTeamAndCheckUser(app.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		handleError(rw, errors.ErrBadRequest)
		return
	}
	// It is not allowed to change the client id yet.
	app.ClientId = mux.Vars(r)["client_id"]

	err = app.Update()
	if err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, app)
}

func (api *Api) appDelete(rw http.ResponseWriter, r *http.Request, user *account.User) {
	app, err := account.FindAppByClientId(mux.Vars(r)["client_id"])
	if err != nil {
		handleError(rw, err)
		return
	}

	if err = app.Delete(*user); err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, app)
}

func (api *Api) appInfo(rw http.ResponseWriter, r *http.Request, user *account.User) {
	app, err := account.FindAppByClientId(mux.Vars(r)["client_id"])
	if err != nil {
		handleError(rw, err)
		return
	}

	_, err = findTeamAndCheckUser(app.Team, user)
	if err != nil {
		handleError(rw, err)
		return
	}

	Ok(rw, app)
}
