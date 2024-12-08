package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type ConnectionConfig struct {
	Host   string
	Port   string
	DbName string
	User   string
	Pass   string
}

func NewDbConnect(c ConnectionConfig) *sql.DB {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", c.Host, c.Port, c.User, c.DbName, c.Pass)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
