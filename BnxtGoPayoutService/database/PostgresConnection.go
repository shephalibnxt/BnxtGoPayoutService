package database

import (
	"database/sql"
	"event-service/config"
	"fmt"

	_ "github.com/lib/pq"
)

func PostgresConnection() *sql.DB {

	//read database configuration from config file
	configInstance := config.InitConfig()
	host := configInstance.String("host")
	port := configInstance.String("port")
	user := configInstance.String("user")
	password := configInstance.String("password")
	dbname := configInstance.String("dbname")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("failed to connect postgres connection", err)

	}

	fmt.Println("Connected to the database successfully..!")
	return db
}
