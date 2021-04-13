package utilities

import (
	Argon "github.com/alexedwards/argon2id"
)

func MakeHash(value string) (string, error) {
	hash, hashError := Argon.CreateHash(value, Argon.DefaultParams)
	if hashError != nil {
		return "", hashError
	}
	return hash, nil
}

func CompareHashes(plaintext string, hash string) (bool, error) {
	match, comparisonError := Argon.ComparePasswordAndHash(plaintext, hash)
	if comparisonError != nil {
		return false, comparisonError
	}
	return match, nil
}
