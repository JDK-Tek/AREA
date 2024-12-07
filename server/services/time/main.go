package main

import (
	// "database/sql"
	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "net/url"
	// "os"
	// "strconv"
	// "time"

	"github.com/gorilla/mux"
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

func timeIn(w http.ResponseWriter, req *http.Request) {
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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/in", timeIn).Methods("POST")
	log.Fatal(http.ListenAndServe(":80", router))
	fmt.Println("time microservice container is running !")
}