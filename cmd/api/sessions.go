package main

import (
	"errors"
	"fmt"
	"github.com/Skaifai/gophers-online-store/internal/data"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Cart struct {
	ID           int64            `json:"id"`
	UserID       int64            `json:"user_id"`
	Total        float64          `json:"total"`
	CartsItems   []*data.CartItem `json:"items"`
	CreationDate time.Time        `json:"creation_date"`
}

func (app *application) recipeGenerateHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	session, err := app.models.ShoppingSessions.Get(id)
	if err != nil {
		return
	}
	items, err := app.models.ShoppingSessions.GetCartItems(id)
	if err != nil {
		return
	}
	user, err := app.models.Users.GetById(session.UserID)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	var itemsString string
	for index, element := range items {
		product, err := app.models.Products.Get(element.ProductID)
		if err != nil {
			return
		}
		itemsString += strconv.Itoa(index+1) + ". " + product.Name +
			" => Quantity: " + strconv.Itoa(int(element.Quantity)) + ", " +
			"Price: " + strconv.FormatFloat(product.Price, 'f', 2, 64) + " " +
			"Total: " + strconv.FormatFloat(float64(element.Quantity)*product.Price, 'f', 2, 64) + "\n"
	}
	fmt.Println(itemsString)

	text := fmt.Sprintf("Cart Id: %d \nBuyer: %s %s \nTotal: %.2f\n%s", session.ID, user.Name, user.Surname, session.Total, itemsString)
	// Set the file name header
	w.Header().Set("Content-Disposition", "attachment; filename=\"recipe.txt\"")
	_, err = fmt.Fprint(w, text)
	if err != nil {
		http.Error(w, "Could not write file", http.StatusInternalServerError)
		return
	}
}

func (app *application) listItemsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	session, err := app.models.ShoppingSessions.Get(id)
	if err != nil {
		return
	}
	items, err := app.models.ShoppingSessions.GetCartItems(id)
	if err != nil {
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"cart": Cart{
		ID:           session.ID,
		UserID:       session.UserID,
		Total:        session.Total,
		CartsItems:   items,
		CreationDate: session.CreationDate,
	}}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addItemToSessionHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SessionID int64 `json:"session_id"`
		ProductID int64 `json:"product_id"`
		Quantity  int64 `json:"quantity"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	item := &data.CartItem{
		SessionID: input.SessionID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}

	err = app.models.CartItems.Insert(item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	total, err := app.models.ShoppingSessions.GetCartTotal(input.SessionID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	session, err := app.models.ShoppingSessions.Get(input.SessionID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	session.Total = total

	err = app.models.ShoppingSessions.Update(session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"cart_item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	item, err := app.models.CartItems.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"cart_item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) removeItemFromSessionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	item, err := app.models.CartItems.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	session, err := app.models.ShoppingSessions.Get(item.SessionID)
	if err != nil {
		return
	}

	err = app.models.CartItems.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	product, err := app.models.Products.Get(item.ProductID)
	if err != nil {
		return
	}
	total := session.Total
	session.Total = total - product.Price*float64(item.Quantity)
	err = app.models.ShoppingSessions.Update(session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "cart item successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateItemInSessionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	item, err := app.models.CartItems.Get(id)
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
		Quantity *int64 `json:"quantity"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if *input.Quantity == 0 {
		err := app.models.CartItems.Delete(id)
		if err != nil {
			return
		}
		err = app.writeJSON(w, http.StatusOK, envelope{"message": "cart_item was deleted due ot quantity is zero"}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
	delta := math.Abs(float64(item.Quantity - *input.Quantity))
	if input.Quantity != nil {
		item.Quantity = *input.Quantity
	}

	err = app.models.CartItems.Update(item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	session, err := app.models.ShoppingSessions.Get(item.SessionID)
	if err != nil {
		return
	}
	product, err := app.models.Products.Get(item.ProductID)
	if err != nil {
		return
	}
	session.Total = session.Total - delta*product.Price
	err = app.models.ShoppingSessions.Update(session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cart_item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
