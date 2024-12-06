package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const API_SEND = "https://discord.com/api/channels/"
const PERMISSIONS = 8 // 2080

type Objects struct {
	Channel int `json:"channel"`
	Message string `json:"message"`
}

type Content struct {
	Dishes Objects `json:"dishes"`
}

func getOAUTHLink(w http.ResponseWriter, req *http.Request) {
	str := "https://discord.com/oauth2/authorize?"
	redirect := req.URL.Query().Get("redirect")
	if redirect == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{ \"error\": \"missing\" }\n")
		return
	}
	x := url.QueryEscape(redirect)
	str += "client_id=" + os.Getenv("DISCORD_ID")
	str += "&permissions=" + strconv.Itoa(PERMISSIONS)
	str += "&response_type=code"
	str += "&redirect_uri=" + x
	str += "&integration_type=0"
	str += "&scope=identify%20email%20bot%20guilds"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, str)
}

func doSomeSend(w http.ResponseWriter, req *http.Request) {
	var content Content

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"token is missing\" }\n")
		return
	}
	data := make(map[string]string)
	data["content"] = content.Dishes.Message
	dataBytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	channel := strconv.Itoa(content.Dishes.Channel)
	fmt.Println(channel, data["content"])
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
	router.HandleFunc("/send", doSomeSend).Methods("POST")
	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}