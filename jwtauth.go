package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/oayomide/xoxo/model"
)

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// unAuths are paths we dont want to authorize for users to access.
		unAuths := []string{"/", "/about", "/blog", "/signup", "/login"}
		var res model.Response

		for _, path := range unAuths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}
		// tokenString := r.Header.Get("Authorization")
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				res.Error = "Unauthorized... no jwt found in cookie"
				json.NewEncoder(w).Encode(res)
				return
			}

			res.Error = "bad request"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}
		tokenString := c.Value
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("xoxo"), nil
		})

		if err != nil {
			// fmt.Println("ERROR AUTHENTICATING JWT WITH MIDDLEWARE\n")
			fmt.Printf("JWT ERROR IS: %v", err)
			// w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusUnauthorized)
			// res.Error = "UNAUTHORIZED... PLEASE CHECK IF TOKEN HAS BEEN PASSED"
			// json.NewEncoder(w).Encode(res)
			// return
			// here, we check if the period the token expired is within 30 seconds of
			// the set expiration time. if not, return bad request
			if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
				res.Error = "bad request... this is where they login in again"
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(res)
				return
			}

			// create new token now
			expirationTime := time.Now().Add(1 * time.Minute)
			claims.ExpiresAt = expirationTime.Unix()
			token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err = token.SignedString([]byte("xoxo"))

			if err != nil {
				res.Error = "internal server error"
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(res)
				return
			}
		}

		if !token.Valid {
			fmt.Printf("OOOOUUUU...ERROR HEEERERE")
			fmt.Printf("ERROR IS:: %s", token)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if token.Valid {
			ctx := context.WithValue(r.Context(), "claims", claims.ID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
	})
}
