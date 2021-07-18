package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Connection string
	DbUser     string
	DbPssword  string
	DbHost     string
	DbPort     string
	DbName     string
}

type Response struct {
	Message string `json:"message"`
}

type User struct {
	Name string `json:"name"`
}

type Rant struct {
	Body       string     `json:"body"`
	User       User       `json:"user"`
	Topics     []string   `json:"topics"`
	QuotedRant QuotedRant `json:"quoted_rant"`
}

type QuotedRant struct {
	Body   string   `json:"body"`
	User   User     `json:"user"`
	Topics []string `json:"topics"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, exist := os.LookupEnv("DB_USER")

	if !exist {
		log.Fatal("DB_USER not set in .env")
	}

	pass, exist := os.LookupEnv("DB_PASS")

	if !exist {
		log.Fatal("DB_PASS not set in .env")
	}

	host, exist := os.LookupEnv("DB_HOST")

	if !exist {
		log.Fatal("DB_HOST not set in .env")
	}

	name, exist := os.LookupEnv("DB_NAME")

	if !exist {
		log.Fatal("DB_NAME not set in .env")
	}

	port, exist := os.LookupEnv("DB_PORT")

	if !exist {
		log.Fatal("DB_PORT not set in .env")
	}

	connection, exist := os.LookupEnv("DB_CONNECTION")

	if !exist {
		log.Fatal("DB_CONNECTION not set in .env")
	}

	_, err := connect(Config{
		DbUser:     user,
		DbPssword:  pass,
		DbHost:     host,
		DbPort:     port,
		DbName:     name,
		Connection: connection,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Response{Message: "unable to connect to database"})
		return
	}

	_ = json.NewEncoder(w).Encode(Response{Message: "ok"})
	return
}

func connect(config Config) (*gorm.DB, error) {
	user := config.DbUser

	pass := config.DbPssword

	host := config.DbHost

	name := config.DbName

	port := config.DbPort

	database := config.Connection

	switch database {
	case "mysql":
		db, err := gorm.Open(config.Connection, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name))

		if err != nil {
			return nil, err
		}

		return db, nil
	case "postgres":
		db, err := gorm.Open(config.Connection, fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, pass, name, port))

		if err != nil {
			return nil, err
		}

		return db, nil
	default:
		return nil, errors.New("database not defined")
	}
}
