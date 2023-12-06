package users

import (
	"net/mail"

	"github.com/pedro-git-projects/chatbot-back/internal/validator"
)

type CreateUserDTO struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Name     string   `json:"name"`
	ImageURL string   `json:"imageUrl"`
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
