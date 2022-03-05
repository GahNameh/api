package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GahNameh/api/src/controller"
	"github.com/GahNameh/api/src/utility"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mvrilo/go-redoc"
	"github.com/rs/cors"
)

func configureStaticFiles(router *mux.Router) {
	fs := http.FileServer(http.Dir(utility.GetEnvPath("public")))
	sfs := http.FileServer(http.Dir(utility.GetEnvPath("public/static")))

	router.Handle("/", fs)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", sfs))
	router.HandleFunc("/favicon.ico", favIcon)
}

func favIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, utility.GetEnvPath("public/static/images/favicon.ico"))
}

func configureEndpoints(router *mux.Router) {
	router.HandleFunc("/api/date/now", controller.GetNow)
	router.HandleFunc("/api/date/{year}/{month}", controller.GetMonthByYearAndId)
}

func configureOpenApi(router *mux.Router) {
	docPath := utility.GetEnvPath("openapi.json")
	doc := redoc.Redoc{SpecFile: docPath, SpecPath: "/docs/openapi.json"}
	docHandler := doc.Handler()
	router.Handle("/docs", docHandler)
	router.Handle("/docs/openapi.json", docHandler)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	configureStaticFiles(myRouter)
	configureEndpoints(myRouter)
	configureOpenApi(myRouter)

	corsHandler := cors.AllowAll().Handler(myRouter)
	log.Fatal(http.ListenAndServe(utility.GetPortString(), corsHandler))
}

func main() {
	fmt.Println("Application Started!")
	handleRequests()
}
