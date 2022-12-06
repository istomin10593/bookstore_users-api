package crypto_utils

import "golang.org/x/crypto/bcrypt"

func HashedValue(input string) (string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
