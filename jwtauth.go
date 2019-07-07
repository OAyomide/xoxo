package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/oayomide/xoxo/model"
)

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// unAuths are paths we dont want to authorize for users to access.
		unAuths := []string{"/", "/about", "/blog"}

		for _, path := range unAuths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}
		tokenString := r.Header.Get("Authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				return nil, fmt.Errorf("INVALID JWT PROVIDED")
			}

			return []byte("xoxo"), nil
		})

		if err != nil {
			var res model.Response
			fmt.Println("ERROR AUTHENTICATING JWT WITH MIDDLEWARE\n")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			res.Error = "UNAUTHORIZED... PLEASE CHECK IF TOKEN HAS BEEN PASSED"
			json.NewEncoder(w).Encode(res)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			ctx := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
	})
}
