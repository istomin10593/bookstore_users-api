package users

// data access object

import (
	"fmt"

	"github.com/istomin10593/bookstore_users-api/datasources/mysql/users_db"
	"github.com/istomin10593/bookstore_users-api/utils/errors"
	"github.com/istomin10593/bookstore_users-api/utils/mysql_utils"
)

const (
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, password, status) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser          = "SELECT Id, first_name, last_name, email, date_created FROM users WHERE id=?;"
	quaryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	quaryDeleteUser       = "DELETE FROM users WHERE id=?;"
	quaryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); getErr != nil {
		return mysql_utils.ParseError(getErr)
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, &user.Password, &user.Status)
	if saveErr != nil {
		mysql_utils.ParseError(saveErr)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		mysql_utils.ParseError(saveErr)
	}

	user.Id = userId

	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(quaryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, updErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if updErr != nil {
		mysql_utils.ParseError(updErr)
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(quaryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, delErr := stmt.Exec(user.Id); delErr != nil {
		return mysql_utils.ParseError(delErr)
	}

	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(quaryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}

	return results, nil
}
