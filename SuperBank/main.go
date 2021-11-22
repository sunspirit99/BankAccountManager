package main

import (
	"SuperBank/controllers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql" //Required for MySQL dialect
)

func main() {

	log.Println("Starting the HTTP server on port 8000")

	router := mux.NewRouter().StrictSlash(true)
	initaliseHandlers(router)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func initaliseHandlers(router *mux.Router) {
	router.HandleFunc("/create", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/creates", controllers.CreateAccountFromCSV).Methods("POST")
	router.HandleFunc("/transfer", controllers.AccountTransfer).Methods("PUT")
	router.HandleFunc("/get", controllers.GetAllAccount).Methods("GET")
	router.HandleFunc("/get/{id}", controllers.GetAccountByID).Methods("GET")
	router.HandleFunc("/update", controllers.UpdateAccountByID).Methods("PUT")
	router.HandleFunc("/delete/{id}", controllers.DeleteAccountByID).Methods("DELETE")
	router.HandleFunc("/delete", controllers.DeleteAccountByID).Methods("DELETE")
	router.HandleFunc("/withdraw", controllers.AccountWithdraw).Methods("PUT")
	router.HandleFunc("/deposit", controllers.AccountDeposit).Methods("PUT")
	router.HandleFunc("/transfer", controllers.AccountTransfer).Methods("PUT")
	router.HandleFunc("/transfers", controllers.AccountTransferFromCSV).Methods("PUT")

}
