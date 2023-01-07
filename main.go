package main

import (
	transaction "Postgres/transactions"
	user2 "Postgres/user"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
	"os"
)

func main() {
	var (
		err error
		db  *gorm.DB
	)

	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	dbpassword := os.Getenv("PASSWORD")

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, dbpassword, dbPort)

	// Opening connection to database
	db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to DB")

	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	db.AutoMigrate(&user2.User{})
	db.AutoMigrate(&transaction.Transaction{})

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

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}

}
