package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GahNameh/api/src/service"
	"github.com/gorilla/mux"
	ptime "github.com/yaa110/go-persian-calendar"
)

func GetNow(w http.ResponseWriter, r *http.Request) {
	pt := ptime.Now()
	response := service.CreateMonthResponse(pt)
	json.NewEncoder(w).Encode(response)
}

func GetMonthByYearAndId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])
	pt := ptime.Date(year, ptime.Month(month), 1, 0, 0, 0, 0, ptime.Iran())
	response := service.CreateMonthResponse(pt)
	json.NewEncoder(w).Encode(response)
}
