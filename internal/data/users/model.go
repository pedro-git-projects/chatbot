package users

import (
	"database/sql"
)

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (email, password, name, role,image_url)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at 
	`

	args := []any{user.Email, user.Password, user.Name, user.Role, user.ImageURL}
	return m.DB.QueryRow(query, args...).Scan(&user.ID, &user.CreatedAt)
}

func (m UserModel) Get(id int64) error {
	return nil
}

func (m UserModel) Update(user *User) error {
	return nil
}

func (m UserModel) Delete(id int64) error {
	return nil
}
