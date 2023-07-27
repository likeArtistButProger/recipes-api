package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"test-api/config"
)

var Conn *sql.DB

func init() {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASS, config.DB_NAME)

	Conn, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = Conn.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Printf("DB successfully connected (%s:%s)!\n", config.DB_HOST, config.DB_PORT)
}
