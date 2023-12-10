package data

import (
	"database/sql"
	"errors"

	"github.com/pedro-git-projects/chatbot-back/src/data/users"
)

var ErrRecordNotFound = errors.New("Registro não encontrado")

type Models struct {
	Users users.UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: users.UserModel{DB: db},
	}
}
