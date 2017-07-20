package main

import (
	"log"
	"net/http"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/gorilla/mux"
	"strconv"
	"encoding/json"
)

const (
	paramName     = "name"
	paramClass    = "class"
	paramSubclass = "var"
	paramPage     = "page"
)

var (
	serverPort = kingpin.Flag("port", "Port to listen on").Short('p').Default("7600").Int()
	staticDir = kingpin.Flag("www-root", "Dir to serve from").Short('d').Default("ui/dist/").String()
)

// Pointer to our asset service (db accessor & cache layer),
// we shouldn't be touching the DB or cache layers at this http layer.
var service *AssetService

// handler that returns all available collection names
func handleSuggestAll(w http.ResponseWriter, r *http.Request) {
	results, err := service.SuggestInitial()
	if err != nil {
		handleError(w, r, err)
		return
	}

	data, err := json.Marshal(results)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(data)
}

// handler for when page & name params are given as part of url params
func handleSuggestName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars[paramName]

	page, err := strconv.Atoi(vars[paramPage])
	if err != nil {
		page = 0
	}

	handleSuggest(w, r, name, "", "", page)
}

// handler for when page, name & class params are given as part of url params
func handleSuggestClass(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars[paramName]
	class := vars[paramClass]

	page, err := strconv.Atoi(vars[paramPage])
	if err != nil {
		page = 0
	}

	handleSuggest(w, r, name, class, "", page)
}

// handler for when page, name, class & subclass params are given as part of url params
func handleSuggestSubclass(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars[paramName]
	class := vars[paramClass]
	subclass := vars[paramSubclass]

	page, err := strconv.Atoi(vars[paramPage])
	if err != nil {
		page = 0
	}

	handleSuggest(w, r, name, class, subclass, page)
}

// handler that writes error to response writer
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(err.Error()))
}

// generic suggest handler that handleSuggest* funcs call
func handleSuggest(w http.ResponseWriter, r *http.Request, name, class, subclass string, page int) {
	matches, err := service.Matches(
		scrub(name),
		scrub(class),
		scrub(subclass),
		page,
	)
	if err != nil {
		handleError(w, r, err)
		return
	}

	data, err := json.Marshal(matches)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(data)
}

// handler for asset fetch endpoint, given name, class & subclass params
func handleFetchAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars[paramName]
	class := vars[paramClass]
	subclass := vars[paramSubclass]

	assetData, err := service.AssetData(
		scrub(name),
		scrub(class),
		scrub(subclass),
	)
	if err != nil {
		handleError(w, r, err)
		return
	}

	data, err := json.Marshal(assetData)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(data)
}

func main() {
	// Parse cli args
	kingpin.Parse()

	// Ready our data & caching service
	svc, err := NewAssetService(NewWysteriaDB())
	if err != nil {
		log.Fatalln(err)
	}
	service = svc
	defer service.Close()

	// Work out what are routes are going to look like
	apiRouteFetch := "/api/1/fetch"
	apiRouteFetchAsset := fmt.Sprintf(
		"%s/{%s}/{%s}/{%s}", apiRouteFetch, paramName, paramClass, paramSubclass,
	)

	apiRouteSuggest := "/api/1/suggest"
	apiRouteSuggestName := fmt.Sprintf(
		"%s/{%s:[0-9]+}/{%s}", apiRouteSuggest, paramPage, paramName,
	)
	apiRouteSuggestClass := fmt.Sprintf(
		"%s/{%s:[0-9]+}/{%s}/{%s}", apiRouteSuggest, paramPage, paramName, paramClass,
	)
	apiRouteSuggestSubclass := fmt.Sprintf(
		"%s/{%s:[0-9]+}/{%s}/{%s}/{%s}", apiRouteSuggest, paramPage, paramName, paramClass, paramSubclass,
	)

	// Set up API routes for fetching data
	router := mux.NewRouter()
	router.HandleFunc(apiRouteFetchAsset, handleFetchAsset).Methods("GET")

	router.HandleFunc(apiRouteSuggest, handleSuggestAll).Methods("GET")
	router.HandleFunc(apiRouteSuggestName, handleSuggestName).Methods("GET")
	router.HandleFunc(apiRouteSuggestClass, handleSuggestClass).Methods("GET")
	router.HandleFunc(apiRouteSuggestSubclass, handleSuggestSubclass).Methods("GET")

	// If it isn't an api route, assume it's looking for a static file
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(*staticDir)))

	addr := fmt.Sprintf(":%d", *serverPort)
	log.Println("Listening On:", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
