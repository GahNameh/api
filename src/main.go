package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	ptime "github.com/yaa110/go-persian-calendar"

	_ "github.com/mattn/go-sqlite3"
	//_ "modernc.org/sqlite"

	"github.com/gorilla/mux"
)

type Response struct {
	Year  int
	Month string
	Days  []Day
}

type Day struct {
	Title     string
	Weekday   string
	IsHoliday bool
	Events    []Event
	IsWeekend bool
	Row       int
	Column    int
}

type Event struct {
	Description string
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
	fmt.Fprintf(w, "Application is running!")
}

func getMonthByYearAndId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])
	pt := ptime.Date(year, ptime.Month(month), 1, 0, 0, 0, 0, ptime.Iran())
	response := generateResponse(pt)
	json.NewEncoder(w).Encode(response)
}

func getNow(w http.ResponseWriter, r *http.Request) {
	pt := ptime.Now()
	response := generateResponse(pt)
	json.NewEncoder(w).Encode(response)
}

func readFromDb(inputYear int, inputMonth int) []DbEvent {
	dbPath := "db.db"
	isHeroku, res := os.LookupEnv("HEROKU")
	if !res && isHeroku == "true" {
		dbPath = "/app/src/" + dbPath
	}
	db, err := sql.Open("sqlite3", dbPath)
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

func generateResponse(pt ptime.Time) Response {
	year := pt.Year()
	month := int(pt.Month())
	events := readFromDb(year, month)

	firstDay := pt.FirstMonthDay()
	lastDay := pt.LastMonthDay()
	firstDayWeekday := int(firstDay.Weekday())

	response := Response{Year: year, Month: pt.Month().String()}
	weekNo := 0
	for i := firstDay.Day(); i <= lastDay.Day(); i++ {

		firstDayWeekday = generateColumn(firstDayWeekday)

		matchedEvents := []DbEvent{}
		for _, event := range events {
			if event.Day == i {
				matchedEvents = append(matchedEvents, event)
			}
		}

		day := Day{
			Title:     strconv.Itoa(i),
			Weekday:   ptime.Weekday(firstDayWeekday).String(),
			IsHoliday: firstDayWeekday == 6,
			IsWeekend: firstDayWeekday == 6,
			Row:       weekNo,
			Column:    firstDayWeekday,
		}
		for _, event := range matchedEvents {
			if event.IsHoliday {
				day.IsHoliday = true
			}
			day.Events = append(day.Events, Event{Description: event.Description})
		}

		response.Days = append(response.Days, day)
		if firstDayWeekday == 6 {
			weekNo++
		}
		firstDayWeekday++
	}
	return response
}

func generateColumn(column int) int {
	if column >= 7 {
		column -= 7
	}
	return column
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/now", getNow)
	myRouter.HandleFunc("/{year}/{month}", getMonthByYearAndId)
	port, res := os.LookupEnv("PORT")
	if !res {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}

func main() {
	fmt.Println("Application Started!")
	handleRequests()
}
