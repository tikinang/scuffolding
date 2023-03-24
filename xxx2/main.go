package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/google/uuid"
	"github.com/tikinang/scuffolding/g"
)

const (
	userId     = "737b7bd2-c5ac-11ed-bf5f-0242ac110002"
	hmacSecret = "abc"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	UserId string `json:"uid,omitempty"`
}

func FreshJwtClaims() JwtClaims {
	now := time.Now()
	return JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "icbaat-server", // FIXME(mpavlicek): Issuer? Server Id?
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * 15)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(), // FIXME(mpavlicek): Replace with database-generated UUID on inserting JWT token to db, if so.
		},
		UserId: userId,
	}
}

// Guest limiting.
// Always limit access on IP.
// For well-behaved clients use Device-Id header, which will be required to be sent with every request on guest access.
// Access limited on IP will be higher than on Device-Id.
// Logged-in users will be limited differently and identified by JWT token.
func main() {
	// login endpoint
	// authorize endpoint
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, FreshJwtClaims())
		signedToken, err := token.SignedString([]byte(hmacSecret))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"jwt_token": signedToken,
		})
	})
	http.HandleFunc("/resource", func(w http.ResponseWriter, r *http.Request) {
		token, err := request.ParseFromRequest(r, g.Empty[request.BearerExtractor](), func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(hmacSecret), nil
		}, request.WithClaims(new(JwtClaims)))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*JwtClaims)
		if !ok {
			fmt.Println("not custom jwt claims")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		exp, err := claims.GetExpirationTime()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if time.Now().After(exp.Time) {
			fmt.Println("expired")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"expire":  exp,
			"user_id": claims.UserId,
		})
	})
	http.ListenAndServe(":8080", http.DefaultServeMux)
}
