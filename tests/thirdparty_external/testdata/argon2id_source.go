package main

import "github.com/alexedwards/argon2id"

func Argon2idCreateHash() string {
	hash, err := argon2id.CreateHash("pa$$word", argon2id.DefaultParams)
	if err != nil {
		return "ERR"
	}
	return hash[:10]
}

func Argon2idCheckPassword() bool {
	hash, _ := argon2id.CreateHash("secret", argon2id.DefaultParams)
	match, err := argon2id.ComparePasswordAndHash("secret", hash)
	if err != nil {
		return false
	}
	return match
}
