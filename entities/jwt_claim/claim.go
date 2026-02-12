package jwt_claim

import "github.com/golang-jwt/jwt/v5"

type Claim struct {
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
