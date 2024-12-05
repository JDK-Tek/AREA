package main

import (
	"fmt"
	"net/http"
	"log"
	// "strconv"

	"github.com/gorilla/mux"
	// "github.com/joho/godotenv"
)

func doSomeSend(w http.ResponseWriter, req *http.Request) {
	// channelStr := req.Header.Get("channel")
	// message := req.Header.Get("message")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello")
}

func main() {
	fmt.Println("discord microservice container is running !")
	router := mux.NewRouter()
	// godotenv.Load("/etc/")
	router.HandleFunc("/send", doSomeSend).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}