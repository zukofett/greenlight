package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/zukofett/greenlight/internal/data"
	"github.com/zukofett/greenlight/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    err := app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, r, err)
        return
    }

    user:= &data.User{
        Name:      input.Name,
        Email:     input.Email,
        Activated: false,
    }

    err = user.Password.Set(input.Password)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    v := validator.New()
    
    if data.ValidateUser(v, user); !v.Valid() {
        app.failedValidationError(w, r, v.Errors)
        return
    }


    ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), app.config.db.queryTimeout)
    defer cancel()

    err = app.models.Users.Insert(ctx, user)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrDuplicateEmail):
            v.AddError("email", "a user with this email address already exists")
            app.failedValidationError(w, r, v.Errors)
        default:
            app.serverErrorResponse(w, r, err)
        }
        return
    }
    
    err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}
