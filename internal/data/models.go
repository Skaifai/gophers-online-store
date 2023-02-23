package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users            UserModel
	Profiles         ProfileModel
	Products         ProductModel
	Comments         CommentModel
	Tokens           TokenModel
	ShoppingSessions ShoppingSessionModel
	CartItems        CartItemModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:            UserModel{DB: db},
		Profiles:         ProfileModel{DB: db},
		Products:         ProductModel{DB: db},
		Comments:         CommentModel{DB: db},
		Tokens:           TokenModel{DB: db},
		ShoppingSessions: ShoppingSessionModel{DB: db},
		CartItems:        CartItemModel{DB: db},
	}
}
