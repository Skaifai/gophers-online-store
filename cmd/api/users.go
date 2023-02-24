package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Skaifai/gophers-online-store/internal/data"
	"github.com/Skaifai/gophers-online-store/internal/validator"
	"github.com/google/uuid"
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

	err = app.models.ShoppingSessions.Insert(user.ID)
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
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
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

	user, err := app.models.Users.GetById(activationLink.UserID)
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

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.ActivationLinks.Update(activationLink)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) authenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if !user.Activated {
		app.notActivatedResponse(w, r)
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

	token, err := data.GenerateTokens(user.ID, user.Username)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if err = app.models.Tokens.SaveToken(token); err != nil {
		app.serverErrorResponse(w, r, err)
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    token.RefreshToken,
		HttpOnly: true,
		MaxAge:   30 * 24 * 60 * 60,
	}

	http.SetCookie(w, &refreshTokenCookie)
	if err = app.writeJSON(w, http.StatusOK, envelope{"refreshToken": token.RefreshToken, "accessToken": token.AccessToken}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) refreshHandler(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		app.UserUnauthorizedResponse(w, r)
	}

	refreshToken := strings.TrimPrefix(authorizationHeader, "Bearer ")
	_, err := data.DecodeRefreshToken(refreshToken)

	if err != nil {
		app.UserUnauthorizedResponse(w, r)
	}

	tokenFromDb, err := app.models.Tokens.FindToken(refreshToken)
	if err != nil {
		app.UserUnauthorizedResponse(w, r)
	}
	userForToken, err := app.models.Users.GetById(tokenFromDb.UserID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	token, err := data.GenerateTokens(userForToken.ID, userForToken.Username)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	if err = app.models.Tokens.SaveToken(token); err != nil {
		app.serverErrorResponse(w, r, err)
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    token.RefreshToken,
		HttpOnly: true,
		MaxAge:   30 * 24 * 60 * 60,
	}

	http.SetCookie(w, &refreshTokenCookie)
	if err = app.writeJSON(w, http.StatusOK, envelope{"refreshToken": token.RefreshToken, "accessToken": token.AccessToken}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		app.UserUnauthorizedResponse(w, r)
	}

	refreshToken := strings.TrimPrefix(authorizationHeader, "Bearer ")

	_, err := data.DecodeRefreshToken(refreshToken)

	if err != nil {
		app.UserUnauthorizedResponse(w, r)
	}

	_, err = app.models.Tokens.FindToken(refreshToken)
	if err != nil {
		app.UserUnauthorizedResponse(w, r)
	}

	err = app.models.Tokens.RemoveToken(refreshToken)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	logoutCookie := http.Cookie{
		Name:   "refreshToken",
		MaxAge: -1,
	}
	http.SetCookie(w, &logoutCookie)
	app.writeJSON(w, http.StatusOK, envelope{"refreshToken": refreshToken}, nil)

}

func (app *application) showUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.models.Users.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.models.Users.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Role        *string    `json:"role"`
		Username    *string    `json:"username"`
		Email       *string    `json:"email"`
		PhoneNumber *string    `json:"phone_number"`
		Password    *string    `json:"password"`
		Name        *string    `json:"name"`
		Surname     *string    `json:"surname"`
		DOB         *time.Time `json:"date_of_birth"`
		Address     *string    `json:"address"`
		AboutMe     *string    `json:"about_me"`
		PictureURL  *string    `json:"picture_url"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Role != nil {
		user.Role = data.RoleType(*input.Role)
	}
	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.PhoneNumber != nil {
		user.PhoneNumber = *input.PhoneNumber
	}
	if input.Password != nil {
		user.Password.Set(*input.Password)
	}
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Surname != nil {
		user.Surname = *input.Surname
	}
	if input.DOB != nil {
		user.DOB = *input.DOB
	}
	if input.Address != nil {
		user.Address = *input.Address
	}
	if input.AboutMe != nil {
		user.AboutMe = *input.AboutMe
	}
	if input.PictureURL != nil {
		user.PictureURL = *input.PictureURL
	}

	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Users.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
