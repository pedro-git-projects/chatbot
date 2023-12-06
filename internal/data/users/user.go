package users

import (
	"time"

	"github.com/pedro-git-projects/chatbot-back/internal/validator"
)

type User struct {
	ID        int64     `json:"id,string"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      UserRole  `json:"role"`
	ImageURL  string    `json:"image_url,omitempty"`
}

func ValidateUser(v *validator.Validator, user *User) {

}
