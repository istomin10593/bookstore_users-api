package users_db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/istomin10593/bookstore_users-api/utils/env"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysql_users_username = "mysql_users_username"
	mysql_users_password = "mysql_users_password"
	mysql_users_host     = "mysql_users_host"
	mysql_users_schema   = "mysql_users_schema"
)

var (
	Client *sql.DB

	username = env.GetEnvVariable(mysql_users_username)
	password = env.GetEnvVariable(mysql_users_password)
	host     = env.GetEnvVariable(mysql_users_host)
	schema   = env.GetEnvVariable(mysql_users_schema)
)

func init() {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username, password, host, schema)

	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}

	log.Println("database successfully configured")
}
