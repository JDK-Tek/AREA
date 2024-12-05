package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"encoding/json"
	"bytes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const API_SEND = "https://discord.com/api/channels/"

func doSomeSend(w http.ResponseWriter, req *http.Request) {
	channel := req.URL.Query().Get("channel")
	token := os.Getenv("DISCORD_TOKEN")
	data := make(map[string]string)
	data["content"] = req.URL.Query().Get("message")
	dataBytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	rep, err := http.NewRequest("POST", API_SEND+channel+"/messages", bytes.NewBuffer(dataBytes))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	rep.Header.Set("Authorization", "Bot " + token)
	rep.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(rep)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{ \"status\": \"%s\" }\n", res.Status)
}

func main() {
	fmt.Println("discord microservice container is running !")
	router := mux.NewRouter()
	godotenv.Load("/usr/mound.d/.env")
	router.HandleFunc("/send", doSomeSend).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}