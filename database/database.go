package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func ConnectDatabase() {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println("Error occured on .env file, please check")
	}
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	user := os.Getenv("DB_USER")
	pass := os.Getenv("PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)
	Database, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("There is an error while connecting to the database ", err)
		panic(err)
	} else {
		Db = Database
		fmt.Println("Successfully connected to DB")
	}
	Db.SetMaxOpenConns(10) // Adjust as needed

	err = Db.Ping()
	if err != nil {
		fmt.Printf("error pinging database: %v \n", err)
	}
	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS stocks (
        	id SERIAL PRIMARY KEY,
			description VARCHAR(255) NOT NULL,
			unit_price INTEGER NOT NULL,
			units INTEGER NOT NULL
    );`)
	if err != nil {
		fmt.Println("Error creating table:", err)

	} else {
		fmt.Println("Table stocks created successfully.")
	}

	_, err = Db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL
	)`)

	if err != nil {
		fmt.Println("Error creating table login:", err)

	} else {
		fmt.Println("Table users created successfully.")
	}

}
