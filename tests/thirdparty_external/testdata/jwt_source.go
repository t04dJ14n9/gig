package main

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

// JWTGetSigningMethod returns the algorithm name of HS256
func JWTGetSigningMethod() string {
	return jwt.SigningMethodHS256.Alg()
}
