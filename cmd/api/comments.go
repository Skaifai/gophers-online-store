package main

import (
	"github.com/Skaifai/gophers-online-store/internal/data"
	"github.com/Skaifai/gophers-online-store/internal/validator"
	"net/http"
	"strings"
)

func (app *application) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProductID    int64  `json:"product_id"`
		CommentOwner int64  `json:"owner_id"`
		Text         string `json:"text"`
	}

	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		app.UserUnauthorizedResponse(w, r)
	}
	// Код ниче почему-то не работает. Токен не находится.
	refreshToken := strings.TrimPrefix(authorizationHeader, "Bearer ")
	//fmt.Println(refreshToken)
	newToken, err := app.models.Tokens.FindToken(refreshToken)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	//fmt.Println(newToken.UserID)
	input.CommentOwner = newToken.UserID
	// Получи юзернейм из этого токена, пожалуйста.

	input.ProductID, err = app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &data.Comment{
		ProductID:    input.ProductID,
		CommentOwner: input.CommentOwner,
		Text:         input.Text,
	}

	v := validator.New()
	if data.ValidateComment(v, comment); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Comments.Insert(comment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"comment": comment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listCommentsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	var input struct {
		ProductID int64
		data.Filters
	}

	input.ProductID = id
	qs := r.URL.Query()
	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "text", "creation_date",
		"-id", "-text", "-creation_date"}

	products, metadata, err := app.models.Comments.GetAll(input.ProductID, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Include the metadata in the response envelope.
	err = app.writeJSON(w, http.StatusOK, envelope{"comments": products, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
