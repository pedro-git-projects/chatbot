package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pedro-git-projects/chatbot-back/src/data/users"
)

type Claims struct {
	UserID int64          `json:"id"`
	Role   users.UserRole `json:"role"`
	jwt.StandardClaims
}

func (app Application) generateJWT(userID int64, role users.UserRole) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(app.config.jwtSecret))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Falha ao assinar token com erro: %v", err))
	}
	return signedToken, nil
}
