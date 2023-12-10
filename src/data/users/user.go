package users

import (
	"time"
)

type User struct {
	ID        int64     `json:"id,string"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Role      UserRole  `json:"role"`
	ImageURL  string    `json:"image_url,omitempty"`
}
