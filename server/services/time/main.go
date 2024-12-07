package main

import (
	"database/sql"
	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "net/url"
	"os"
	// "strconv"
	// "time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const API = "https://tools.aimylogic.com/api/now?tz=Europe/Paris"

type Time struct {
	HowMuch int `json:"howmuch"`
	Unit string `json:"unit"`
}

type Content struct {
	Spices Time `json:"spices"`
}

type Response struct {
	Timestamp int64 `json:"timestamp"`
}

func makeAnError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
}

func timeIn(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var userContent Content
	var timeRep Response

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&userContent)
	if err != nil {
		makeAnError(w, err, http.StatusBadRequest)
		return
	}
	rep, err := http.Get(API)
	if err != nil {
		makeAnError(w, err, http.StatusBadGateway)
		return
	}
	defer rep.Body.Close()
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		makeAnError(w, err, http.StatusBadGateway)
		return
	}
	err = json.Unmarshal(body, &timeRep)
	if err != nil {
		makeAnError(w, err, http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello %v", timeRep.Timestamp)
}

func connectToDatabase() (*sql.DB, error) {
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD not found")
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER not found")
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal("DB_HOST not found")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME not found")
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Fatal("DB_PORT not found")
	}
	connectStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)
	return sql.Open("postgres", connectStr)
}

func miniProxy(f func(http.ResponseWriter, *http.Request, *sql.DB), c *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(a http.ResponseWriter, b *http.Request) {
		f(a, b, c)
	}
}

func main() {
	db, err := connectToDatabase()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(84)
	}
	fmt.Println("time microservice container is running !")
	router := mux.NewRouter()
	router.HandleFunc("/in", miniProxy(timeIn, db)).Methods("POST")
	log.Fatal(http.ListenAndServe(":80", router))
}