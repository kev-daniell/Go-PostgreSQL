package main

import (
	transaction "Postgres/transactions"
	"Postgres/types"
	user2 "Postgres/user"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"

	"net/http"
)

func main() {
	var (
		err error
		db  *gorm.DB
	)

	host := os.Getenv("HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DBPORT"))
	if err != nil {
		panic(err)
	}
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	dbpassword := os.Getenv("PASSWORD")

	// Database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable", host, user, dbpassword, dbname, dbPort)

	// Opening connection to database
	fmt.Println("dialect:", dsn)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to DB")

	err = db.AutoMigrate(&types.User{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&types.Transaction{})
	if err != nil {
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode("Hi Mom")
		if err != nil {
			return
		}
	})

	txnHandler := transaction.NewTxn(db)
	userHandler := user2.NewUser(db)

	router.HandleFunc("/transaction", txnHandler.MakeTransaction).Methods("POST")
	router.HandleFunc("/transaction", txnHandler.GetTransaction).Methods("GET")
	router.HandleFunc("/users", userHandler.MakeUser).Methods("POST")
	router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	router.HandleFunc("/user/{id}", userHandler.GetUser).Methods("Get")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}

}
