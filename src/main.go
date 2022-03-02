package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	//_ "modernc.org/sqlite"

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

type DbEvent struct {
	Id          int
	Year        int
	Month       int
	Day         int
	Type        int
	IsHoliday   bool
	Description string
}

var Result Response

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func getMonth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Result)
}

func getMonthById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Println(key)
	json.NewEncoder(w).Encode(Result)
}

func getMonthByYearAndId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])
	events := readFromDb(year, month)
	json.NewEncoder(w).Encode(events)
}

func readFromDb(inputYear int, inputMonth int) []DbEvent {
	db, err := sql.Open("sqlite3", "/app/src/db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM Holiday WHERE Year = %d AND Month = %d", inputYear, inputMonth))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	entities := make([]DbEvent, 0)

	for rows.Next() {
		entity := DbEvent{}
		err = rows.Scan(&entity.Id,
			&entity.Year,
			&entity.Month,
			&entity.Day,
			&entity.Type,
			&entity.IsHoliday,
			&entity.Description)

		if err != nil {
			log.Fatal(err)
		}
		entities = append(entities, entity)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return entities
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/month", getMonth)
	myRouter.HandleFunc("/month/{id}", getMonthById)
	myRouter.HandleFunc("/{year}/{month}", getMonthByYearAndId)
	port, res := os.LookupEnv("PORT")
	if !res {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}

func main() {
	fmt.Println("Welcome here!")
	listAll()
	mockResponse := []Day{
		{Title: "1"},
		{Title: "2"},
	}
	Result = Response{Days: mockResponse}
	handleRequests()
}

func listAll() {
	fmt.Println("Start")
	fmt.Println("/app/src")
	files, err := ioutil.ReadDir("/app/src/database")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
	fmt.Println("End")
}
