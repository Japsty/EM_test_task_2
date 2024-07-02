package connect

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

// NewPostgresConnection - функция создающая подключение к БД и предоставляющая его наружу
func NewPostgresConnection(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("Error parsing database config", err)
		return nil, err
	}

	time.Sleep(3 * time.Second)

	err = conn.Ping()
	if err != nil {
		fmt.Println("Error pinging the database:", err)
		return nil, err
	}

	fmt.Println("DB connection opened")

	return conn, nil
}
