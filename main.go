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

	// transactions:
	router.HandleFunc("/transactions", txnHandler.MakeTransaction).Methods(http.MethodPost)
	router.HandleFunc("/transactions", txnHandler.GetAllTransactions).Methods(http.MethodGet)
	router.HandleFunc("/transactions/{id}", txnHandler.GetTransaction).Methods(http.MethodGet)
	router.HandleFunc("/transactions/{id}", txnHandler.UpdateTransaction).Methods(http.MethodPatch)
	router.HandleFunc("/transactions/{id}", txnHandler.DeleteTransaction).Methods(http.MethodDelete)

	// users:
	router.HandleFunc("/users", userHandler.MakeUser).Methods(http.MethodPost)
	router.HandleFunc("/users", userHandler.GetAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", userHandler.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods(http.MethodPatch)
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods(http.MethodDelete)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		return
	}
}
