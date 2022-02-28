package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Response struct {
	Days []Day
}

type Day struct {
	Title       string
	Description string
	IsHoliday   bool
	IsWeekend   bool
	Row         int
	Column      int
}

var Result Response

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func getMonth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Println(key)
	json.NewEncoder(w).Encode(Result)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/month/{id}", getMonth)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), myRouter))
}

func main() {
	fmt.Println("Welcome here!")
	mockResponse := []Day{
		{Title: "1"},
		{Title: "2"},
	}
	Result = Response{Days: mockResponse}
	handleRequests()
}
