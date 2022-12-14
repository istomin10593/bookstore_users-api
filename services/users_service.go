package services

import (
	"github.com/istomin10593/bookstore_users-api/domain/users"
	"github.com/istomin10593/bookstore_users-api/utils/crypto_utils"
	"github.com/istomin10593/bookstore_users-api/utils/date"
	"github.com/istomin10593/bookstore_utils-go/logger"
	"github.com/istomin10593/bookstore_utils-go/rest_errors"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct {
}

type usersServiceInterface interface {
	GetUser(int64) (*users.User, rest_errors.RestErr)
	CreateUser(users.User) (*users.User, rest_errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, rest_errors.RestErr)
	DeleteUser(int64) rest_errors.RestErr
	SearchUser(string) (users.Users, rest_errors.RestErr)
	LoginUser(users.LoginRequest) (*users.User, rest_errors.RestErr)
}

func (s *usersService) GetUser(userId int64) (*users.User, rest_errors.RestErr) {
	result := &users.User{Id: userId}

	if err := result.Get(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *usersService) CreateUser(user users.User) (*users.User, rest_errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Status = users.StatusActive
	user.DateCreated = date.GetNowDBFormat()
	password, err := crypto_utils.HashedValue(user.Password)
	if err != nil {
		logger.Error("error when trying to get hashed value", err)
		restErr := rest_errors.NewInternalServerError("database error", err)
		return nil, restErr
	}
	user.Password = password

	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, rest_errors.RestErr) {
	current, err := UsersService.GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}

		if user.LastName != "" {
			current.LastName = user.LastName
		}

		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}

	if err := current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *usersService) DeleteUser(userId int64) rest_errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func (s *usersService) SearchUser(status string) (users.Users, rest_errors.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

func (s *usersService) LoginUser(request users.LoginRequest) (*users.User, rest_errors.RestErr) {
	password, err := crypto_utils.HashedValue(request.Password)
	if err != nil {
		logger.Error("error when trying to get hashed value", err)
		restErr := rest_errors.NewInternalServerError("database error", err)
		return nil, restErr
	}
	dao := &users.User{
		Email:    request.Email,
		Password: password,
	}

	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}

	return dao, nil
}
