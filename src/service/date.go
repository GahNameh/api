package service

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/GahNameh/api/src/entity"
	"github.com/GahNameh/api/src/model"
	"github.com/GahNameh/api/src/utility"

	ptime "github.com/yaa110/go-persian-calendar"

	_ "github.com/mattn/go-sqlite3"
)

func generateColumn(column int) int {
	if column >= 7 {
		column -= 7
	}
	return column
}

func readFromDb(inputYear int, inputMonth int) []entity.Event {
	db, err := sql.Open("sqlite3", utility.GetEnvPath("db.db"))
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

	entities := make([]entity.Event, 0)

	for rows.Next() {
		entity := entity.Event{}
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

func CreateMonthResponse(pt ptime.Time) model.Response {
	year := pt.Year()
	month := int(pt.Month())
	events := readFromDb(year, month)

	firstDay := pt.FirstMonthDay()
	lastDay := pt.LastMonthDay()
	firstDayWeekday := int(firstDay.Weekday())

	response := model.Response{Year: year, Month: pt.Month().String()}
	weekNo := 0
	for i := firstDay.Day(); i <= lastDay.Day(); i++ {

		firstDayWeekday = generateColumn(firstDayWeekday)

		matchedEvents := []entity.Event{}
		for _, event := range events {
			if event.Day == i {
				matchedEvents = append(matchedEvents, event)
			}
		}

		day := model.Day{
			Title:     strconv.Itoa(i),
			Value:     fmt.Sprintf("%d/%d/%d", year, month, i),
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
			day.Events = append(day.Events, model.Event{Description: event.Description})
		}

		response.Days = append(response.Days, day)
		if firstDayWeekday == 6 {
			weekNo++
		}
		firstDayWeekday++
	}
	return response
}
