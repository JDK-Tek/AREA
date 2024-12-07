package main

import (
	"database/sql"
	"errors"
	"time"

	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	// "net/url"
	"os"
	// "strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const API = "https://tools.aimylogic.com/api/now?tz=Europe/Paris"

type Spices struct {
	HowMuch int `json:"howmuch"`
	Unit string `json:"unit"`
}

type Content struct {
	Spices Spices `json:"spices"`
	BridgeId int `json:"bridge"`
}

type Response struct {
	Timestamp int64 `json:"timestamp"`
}

func makeAnError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
}

func getTimeNow(w http.ResponseWriter, req *http.Request) (c Content, n int64, err error) {
	var r Response

	c.BridgeId = -1
	c.Spices.HowMuch = -1
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&c)
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
	err = json.Unmarshal(body, &r)
	if err != nil {
		makeAnError(w, err, http.StatusBadGateway)
		return
	}
	if c.BridgeId == -1 || c.Spices.HowMuch == -1 {
		err = errors.New("invalid parsing")
		makeAnError(w, err, http.StatusBadRequest)
		return
	}
	n = r.Timestamp
	return
}

func spices2Seconds(spices Spices) int64 {
	var association = map[string]int{
		"weeks": 7 * 24 * 3600,
		"days": 24 * 3600,
		"hours": 3600,
		"minutes": 60,
		"seconds": 1,
	}
	return int64(spices.HowMuch * association[spices.Unit])
}

func timeIn(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var maybeid int

	content, ms, err := getTimeNow(w, req)
	if err != nil {
		return
	}
	err = db.QueryRow("select id from micro_time where bridgeid = $1", content.BridgeId).Scan(&maybeid)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{ \"error\": \"it already exists\" }\n")
		return
	}
	secs := ms / 1000
	nsecs := (ms % 1000) * 1e6
	secs += spices2Seconds(content.Spices)
	timestamp := time.Unix(secs, nsecs)
	_, err = db.Exec("insert into micro_time (bridgeid, triggers) values ($1, $2)", content.BridgeId, timestamp)
	fmt.Println("hello man !")
	if err != nil {
		makeAnError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello %v", timestamp)
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
	// connectStr := fmt.Sprintf(
	// 	"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
	// 	dbUser,
	// 	dbPassword,
	// 	dbHost,
	// 	dbPort,
	// 	dbName,
	// )
	connectStr := fmt.Sprintf(
		"postgresql://%s:%s@database:5432/area_database?sslmode=disable",
		dbUser,
		dbPassword,
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