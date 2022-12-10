package users

// data access object

import (
	"fmt"
	"strings"

	"github.com/istomin10593/bookstore_users-api/datasources/mysql/users_db"
	"github.com/istomin10593/bookstore_utils-go/logger"
	"github.com/istomin10593/bookstore_utils-go/rest_errors"
)

const (
	queryInsertUser             = "INSERT INTO users(first_name, last_name, email, date_created, password, status) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	quaryUpdateUser             = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	quaryDeleteUser             = "DELETE FROM users WHERE id=?;"
	quaryFindByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	quaryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status=?;"

	ErrorNoRows = "no rows in result set"
)

func (user *User) Get() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("error when trying to get user", rest_errors.NewRestError("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if getErr := result.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.DateCreated,
		&user.Status); getErr != nil {
		logger.Error("error when trying to get user by id", getErr)
		return rest_errors.NewInternalServerError("error when trying to get user", rest_errors.NewRestError("database error"))
	}

	return nil
}

func (user *User) Save() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("error when trying to save user", rest_errors.NewRestError("database error"))
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(
		user.FirstName,
		user.LastName,
		user.Email,
		user.DateCreated,
		user.Password,
		user.Status)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return rest_errors.NewInternalServerError("error when trying to save user", rest_errors.NewRestError("database error"))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get the last insert id after creating a new user", err)
		return rest_errors.NewInternalServerError("error when trying to save user", rest_errors.NewRestError("database error"))
	}

	user.Id = userId

	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(quaryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("error when trying to update user", rest_errors.NewRestError("database error"))
	}
	defer stmt.Close()

	_, updErr := stmt.Exec(
		user.FirstName,
		user.LastName,
		user.Email,
		user.Id)
	if updErr != nil {
		logger.Error("error when trying to update user", updErr)
		return rest_errors.NewInternalServerError("error when trying to update user", rest_errors.NewRestError("database error"))
	}

	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(quaryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", rest_errors.NewRestError("database error"))
	}
	defer stmt.Close()

	if _, delErr := stmt.Exec(user.Id); delErr != nil {
		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", rest_errors.NewRestError("database error"))
	}

	return nil
}

func (user *User) FindByStatus(status string) ([]User, *rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(quaryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find users by status", rest_errors.NewRestError("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find users by status", rest_errors.NewRestError("database error"))
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.DateCreated,
			&user.Status); err != nil {
			logger.Error("error when trying to scan user row into user struct", err)
			return nil, rest_errors.NewInternalServerError("error when trying to find users by status", rest_errors.NewRestError("database error"))
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status), nil)
	}

	return results, nil
}

func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(quaryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("error when trying to find users by email and password", rest_errors.NewRestError("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)

	if getErr := result.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.DateCreated,
		&user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), ErrorNoRows) {
			return rest_errors.NewNotFoundError("invalid user credentials", getErr)
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return rest_errors.NewInternalServerError("error when trying to find users by email and password", rest_errors.NewRestError("database error"))
	}

	return nil
}
