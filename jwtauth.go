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

		if !token.Valid {
			fmt.Println("JWT TOKEN IS NOT VALID. HERE IS WHERE WE REFRESH IF ITS BEEN WITHIN 30s WHEN THE TOKEN EXPIRED")
			if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
				res.Error = err.Error()
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(res)
				return
			} else {
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

				ctx := context.WithValue(r.Context(), "claims", claims.ID.Hex())
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			}
			// res.Error = err.Error()
			// w.WriteHeader(http.StatusUnauthorized)
			// json.NewEncoder(w).Encode(res)
			// return
		}

		if err != nil {
			fmt.Printf("JWT ERROR IS: %v", err)
			if err == jwt.ErrSignatureInvalid {
				fmt.Println("UNAUTHORIZED. JWT SIG IS INVALID")
				res.Error = err.Error()
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(res)
				return
			}

			// w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusUnauthorized)
			// res.Error = "UNAUTHORIZED... PLEASE CHECK IF TOKEN HAS BEEN PASSED"
			// json.NewEncoder(w).Encode(res)
			// return
			// here, we check if the period the token expired is within 30 seconds of
			// the set expiration time. if not, return bad request
		}
		if token.Valid {
			ctx := context.WithValue(r.Context(), "claims", claims.ID.Hex())
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
	})
}
