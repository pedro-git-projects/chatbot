package users

import (
	"net/mail"

	"github.com/pedro-git-projects/chatbot-back/internal/validator"
)

type CreateUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

func (dto CreateUserDTO) Validate(v *validator.Validator) {
	v.Check(dto.Name != "", "name", "é obrigatório")
	v.Check(dto.Password != "", "senha", "é obrigatória")
	v.Check(dto.Email != "", "senha", "é obrigatório")

	_, err := mail.ParseAddress(dto.Email)
	validMail := err == nil
	v.Check(validMail, "email", "deve ser um endereço de email válido")
}
