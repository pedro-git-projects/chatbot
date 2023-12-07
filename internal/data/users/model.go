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

	user := User{}
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

func (m UserModel) Update(id int64, updatedUser *User) (*User, error) {
	existingUser, err := m.Get(id)
	if err != nil {
		return nil, err
	}

	if updatedUser.Email != "" {
		existingUser.Email = updatedUser.Email
	}
	if updatedUser.Password != "" {
		existingUser.Password = updatedUser.Password
	}
	if updatedUser.Name != "" {
		existingUser.Name = updatedUser.Name
	}
	if updatedUser.Role != "" {
		existingUser.Role = updatedUser.Role
	}
	if updatedUser.ImageURL != "" {
		existingUser.ImageURL = updatedUser.ImageURL
	}

	query := `
		UPDATE users
		SET email = $1, password = $2, name = $3, role = $4, image_url = $5
		WHERE id = $6
		RETURNING id, email, name, role, image_url, created_at
	`

	err = m.DB.QueryRow(query, existingUser.Email, existingUser.Password, existingUser.Name, existingUser.Role, existingUser.ImageURL, id).
		Scan(
			&existingUser.ID,
			&existingUser.Email,
			&existingUser.Name,
			&existingUser.Role,
			&existingUser.ImageURL,
			&existingUser.CreatedAt,
		)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Atualização falhou com erro: %v", err))
	}

	return existingUser, nil
}

func (m UserModel) Delete(id int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return errors.New(fmt.Sprintf("Falha ao excluir usuário: %v", err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New(fmt.Sprintf("Falha ao obter número de linhas afetadas: %v", err))
	}

	if rowsAffected == 0 {
		return errors.New("Usuário não encontrado")
	}

	return nil
}
