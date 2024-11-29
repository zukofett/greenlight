package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/zukofett/greenlight/internal/data"
	"github.com/zukofett/greenlight/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationError(w, r, v.Errors)
		return
	}

	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), app.config.db.queryTimeout)
	defer cancel()

	user, err := app.models.Users.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	ctx, cancel = context.WithTimeout(context.WithoutCancel(r.Context()), app.config.db.queryTimeout)
	defer cancel()

	token, err := app.models.Tokens.New(ctx, user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
