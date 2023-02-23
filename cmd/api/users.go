package main

import (
	"errors"
	"github.com/Skaifai/gophers-online-store/internal/data"
	"github.com/Skaifai/gophers-online-store/internal/validator"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string    `json:"name"`
		Surname     string    `json:"surname"`
		Username    string    `json:"username"`
		DOB         time.Time `json:"date_of_birth"`
		PhoneNumber string    `json:"phone_number"`
		Address     string    `json:"address"`
		Email       string    `json:"email"`
		Password    string    `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:        input.Name,
		Surname:     input.Surname,
		Username:    input.Username,
		DOB:         input.DOB,
		PhoneNumber: input.PhoneNumber,
		Address:     input.Address,
		Email:       input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	uuidCode := strings.Replace(uuid.New().String(), "-", "", -1)
	err = app.models.ActivationLinks.Insert(user, uuidCode)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"name":  user.Name,
			"email": user.Email,
			"uuid":  uuidCode,
		}
		err = app.mailer.Send(user.Email, "user.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	uuidParam, err := app.readUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	activationLink, err := app.models.ActivationLinks.Get(uuidParam)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user, err := app.models.Users.Get(activationLink.UserID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if uuidParam == activationLink.Link {
		user.Activated = true
		activationLink.Activated = true
	} else {
		app.badRequestResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
