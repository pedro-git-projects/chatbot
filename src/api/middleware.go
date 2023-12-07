package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func (app application) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedResponse(w, r, "Faltando cabeçalhos de autorização")
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			app.unauthorizedResponse(w, r, "Cabeçalho de autorização malformado")
			return
		}

		tokenString := tokenParts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Método de assinatura ineseprada: %v", token.Header["alg"])
			}

			return []byte(app.config.jwtSecret), nil
		})

		if err != nil {
			app.unauthorizedResponse(w, r, "Token inválido: "+err.Error())
			return
		}

		if !token.Valid {
			app.unauthorizedResponse(w, r, "Token inválido")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			app.unauthorizedResponse(w, r, "Alegações de token inválidas")
			return
		}

		userID, ok := claims["id"].(float64)
		if !ok {
			app.unauthorizedResponse(w, r, "ID de usuário inválido nas alegações do token")
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			app.unauthorizedResponse(w, r, "Papel de usuário inválido nas alegações do token")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", int64(userID))
		ctx = context.WithValue(ctx, "role", role)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
