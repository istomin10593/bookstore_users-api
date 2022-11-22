package services

import (
	"github.com/istomin10593/bookstore_users-api/domain/users"
	"github.com/istomin10593/bookstore_users-api/utils/errors"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	return &user, nil

}
