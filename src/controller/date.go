package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GahNameh/api/src/model"
	"github.com/GahNameh/api/src/service"
	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
	ptime "github.com/yaa110/go-persian-calendar"
)

func getRequestFromBody(r *http.Request) model.Request {
	var request model.Request
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			return request
		}

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&request)
		if err != nil {
			return request
		}
	}
	return request
}

func getRequestFromQuery(r *http.Request) model.Request {
	var request model.Request
	query := r.URL.Query()
	requestFormat := query["format"]
	if requestFormat != nil {
		request.Format = requestFormat[0]
	}
	return request
}

func GetNow(w http.ResponseWriter, r *http.Request) {
	request := getRequestFromQuery(r)
	pt := ptime.Now()
	response := service.CreateMonthResponse(pt, request)
	json.NewEncoder(w).Encode(response)
}

func GetMonthByYearAndId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])

	request := getRequestFromQuery(r)

	pt := ptime.Date(year, ptime.Month(month), 1, 0, 0, 0, 0, ptime.Iran())
	response := service.CreateMonthResponse(pt, request)
	json.NewEncoder(w).Encode(response)
}
