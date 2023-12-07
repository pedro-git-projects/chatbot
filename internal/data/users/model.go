package users

import (
	"database/sql"
	"errors"
	"fmt"
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

func (m UserModel) Authenticate(email, password string) (*User, error) {
	query := `
	SELECT id, email, name, role, image_url, created_at
	FROM users
	WHERE email = $1 AND password = $2
	`

	user := User{}
	err := m.DB.QueryRow(query, email, password).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.ImageURL,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("Usuário não encontrado")
	} else if err != nil {
		return nil, errors.New(fmt.Sprintf("Query falhou com erro: %v", err))
	}

	return &user, nil
}

func (m UserModel) Get(id int64) (*User, error) {
	query := `
		SELECT id, email, name, role, image_url, created_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := m.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.ImageURL,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("Usuário não encontrado")
	} else if err != nil {
		return nil, errors.New(fmt.Sprintf("Query falhou com erro: %v", err))
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	return nil
}

func (m UserModel) Delete(id int64) error {
	return nil
}
