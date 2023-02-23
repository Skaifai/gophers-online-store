package main

import (
	"errors"
	"github.com/Skaifai/gophers-online-store/internal/data"
	"github.com/Skaifai/gophers-online-store/internal/validator"
	"net/http"
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

	//token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//	return
	//}
	//app.background(func() {
	//	data := map[string]any{
	//		"activationToken": token.Plaintext,
	//		"userID":          user.ID,
	//	}
	//	err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
	//	if err != nil {
	//		app.logger.PrintError(err, nil)
	//	}
	//})
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
