package service

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/GahNameh/api/src/entity"
	"github.com/GahNameh/api/src/model"
	"github.com/GahNameh/api/src/utility"

	gcache "github.com/patrickmn/go-cache"
	ptime "github.com/yaa110/go-persian-calendar"

	_ "github.com/mattn/go-sqlite3"
)

var c = gcache.New(5*time.Minute, 10*time.Minute)

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

func generateResponse(pt ptime.Time, format string) model.Response {
	var response model.Response
	key := fmt.Sprintf("%v%v", pt.Year(), int(pt.Month()))
	cached, isFound := c.Get(key)
	if isFound {
		cachedResponse := cached.(*model.Response)
		response = *cachedResponse
	} else {
		response = generateMonth(pt)
		c.Set(key, &response, gcache.NoExpiration)
	}
	generateResponseValues(&response, format)
	return response
}

func generateMonth(pt ptime.Time) model.Response {
	year := pt.Year()
	month := int(pt.Month())
	events := readFromDb(year, month)

	firstDay := pt.FirstMonthDay()
	lastDay := pt.LastMonthDay()
	firstDayWeekday := int(firstDay.Weekday())

	response := model.Response{Year: year, Month: pt.Month().String(), MonthId: month}

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
			Title: strconv.Itoa(i),
			//Value:     fmt.Sprintf("%d/%d/%d", year, month, i),
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

func generateResponseValues(response *model.Response, format string) {
	format = strings.Replace(format, "y", strconv.Itoa(response.Year), 1)
	if strings.Contains(format, "mm") {
		format = strings.Replace(format, "mm", fmt.Sprintf("%02d", response.MonthId), 1)
	} else {
		format = strings.Replace(format, "m", strconv.Itoa(response.MonthId), 1)
	}
	for i, day := range response.Days {
		if strings.Contains(format, "dd") {
			response.Days[i].Value = strings.Replace(format, "dd", fmt.Sprintf("%02s", day.Title), 1)
		} else {
			response.Days[i].Value = strings.Replace(format, "d", day.Title, 1)
		}
	}
}

func getRequestFormat(request model.Request) string {
	if request.Format == "" {
		return "y/m/d"
	}
	return request.Format
}

func CreateMonthResponse(pt ptime.Time, request model.Request) model.Response {
	requestFormat := getRequestFormat(request)
	return generateResponse(pt, requestFormat)
}
