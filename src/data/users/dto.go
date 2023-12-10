package users

import (
	"net/mail"

	"github.com/pedro-git-projects/chatbot-back/src/data/validator"
)

type CreateUserDTO struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Name     string   `json:"name"`
	ImageURL string   `json:"image_url"`
	Role     UserRole `json:"role"`
}

func (dto CreateUserDTO) Validate(v *validator.Validator) {
	v.Check(dto.Name != "", "name", "é obrigatório")
	v.Check(dto.Password != "", "senha", "é obrigatória")
	v.Check(dto.Email != "", "senha", "é obrigatório")

	validRole := dto.Role == "user" || dto.Role == "admin" || dto.Role == "collaborator"
	v.Check(validRole, "role", "deve ser uma das opções (admin|collaborator|user)")

	_, err := mail.ParseAddress(dto.Email)
	validMail := err == nil
	v.Check(validMail, "email", "deve ser um endereço de email válido")
}

type LoginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto LoginUserDTO) Validate(v *validator.Validator) {
	v.Check(dto.Password != "", "senha", "é obrigatória")
	v.Check(dto.Email != "", "senha", "é obrigatório")

	_, err := mail.ParseAddress(dto.Email)
	validMail := err == nil
	v.Check(validMail, "email", "deve ser um endereço de email válido")
}

type UpdateUserDTO struct {
	Email    string   `json:"email,omitempty"`
	Password string   `json:"password,omitempty"`
	Name     string   `json:"name,omitempty"`
	ImageURL string   `json:"imageUrl,omitempty"`
	Role     UserRole `json:"role,omitempty"`
}

func (dto UpdateUserDTO) Validate(v *validator.Validator) {
	if dto.Role != "" {
		validRole := dto.Role == "user" || dto.Role == "admin" || dto.Role == "collaborator"
		v.Check(validRole, "role", "deve ser uma das opções (admin|collaborator|user)")
	}

	if dto.Email != "" {
		_, err := mail.ParseAddress(dto.Email)
		validMail := err == nil
		v.Check(validMail, "email", "deve ser um endereço de email válido")
	}
}
