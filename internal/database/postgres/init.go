package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	var err error

	// Retry logic
	for i := 0; i < 5; i++ {
		connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", c.Host, c.Port, c.User, c.DbName, c.Pass)
		db, err := sql.Open("postgres", connString)
		if err == nil && db.Ping() == nil {
			log.Println("Database connection established!")
			return db
		}

		log.Printf("Retrying database connection in 5 seconds... (%d/5)\n", i+1)
		time.Sleep(5 * time.Second)
	}

	log.Fatalf("Failed to connect to the database after 5 retries: %v", err)
	return nil
}
