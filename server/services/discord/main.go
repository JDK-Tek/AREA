package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	// "io/ioutil"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const API_SEND = "https://discord.com/api/channels/"
const API_OAUTH = "https://discord.com/api/oauth2/token"
const API_USER = "https://discord.com/api/v10/users/@me"

const PERMISSIONS = 8 // 2080

const EXPIRATION = 60 * 30

type Objects struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

type Content struct {
	Dishes Objects `json:"spices"`
}

func getOAUTHLink(w http.ResponseWriter, req *http.Request) {
	str := "https://discord.com/oauth2/authorize?"
	str += "client_id=" + os.Getenv("DISCORD_ID")
	str += "&permissions=" + strconv.Itoa(PERMISSIONS)
	str += "&response_type=code"
	str += "&redirect_uri=" + url.QueryEscape(os.Getenv("DISCORD_REDIRECT"))
	str += "&integration_type=0"
	str += "&scope=identify+email+bot+guilds"
	fmt.Println(str)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, str)
}

type Result struct {
	Code string `json:"code"`
}

type TokenResult struct {
	Token string `json:"access_token"`
}

type UserResult struct {
	ID string `json:"id"`
}

func setOAUTHToken(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var res Result
	var tok TokenResult
	var user UserResult
	var tokid int
	var owner = -1
	// make the request to discord api
	clientid := os.Getenv("DISCORD_ID")
	clientsecret := os.Getenv("DISCORD_SECRET")
	data := url.Values{}
	err := json.NewDecoder(req.Body).Decode(&res)
	if err != nil {
		fmt.Fprintln(w, "decode", err.Error())
		return
	}
	data.Set("client_id", clientid)
	data.Set("client_secret", clientsecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", res.Code)
	data.Set("redirect_uri", os.Getenv("DISCORD_REDIRECT"))
	rep, err := http.PostForm(API_OAUTH, data);
	if err != nil {
		fmt.Fprintln(w, "postform", err.Error())
		return
	}
	defer rep.Body.Close()
	err = json.NewDecoder(rep.Body).Decode(&tok)
	if err != nil {
		fmt.Fprintln(w, "decode", err.Error())
		return
	}

	// make the request for the user
	req, err = http.NewRequest("GET", API_USER, nil)
	if err != nil {
		fmt.Fprintln(w, "request error", err.Error())
		return
	}
	req.Header.Set("Authorization", "Bearer "+tok.Token)
	client := &http.Client{}
	rep, err = client.Do(req)
	if err != nil {
		fmt.Fprintln(w, "client do", err.Error())
		return
	}
	defer rep.Body.Close()
	err = json.NewDecoder(rep.Body).Decode(&user)
	if err != nil {
		fmt.Fprintln(w, "decode", err.Error())
		return
	}

	// seelect the user id shit
	err = db.QueryRow("select id, owner from tokens where userid = $1", user.ID).Scan(&tokid, &owner)
	if err != nil {
		err = db.QueryRow("insert into tokens (service, token, userid) values ($1, $2, $3) returning id",
			"discord",
			tok.Token,
			user.ID,
		).Scan(&tokid)
		if err != nil {
			fmt.Fprintln(w, "db insert", err.Error())
			return
		}
		err = db.QueryRow("insert into users (tokenid) values ($1) returning id", tokid).Scan(&owner)
		if err != nil {
			fmt.Fprintln(w, "db insert", err.Error())
			return
		}
		db.Exec("update tokens set owner = $1 where id = $2", owner, tokid)
	}

	// create the token
	secretBytes := []byte(os.Getenv("BACKEND_KEY"))
	claims := jwt.MapClaims{
		"id":  owner,
		"exp": time.Now().Add(time.Second * EXPIRATION).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secretBytes)
	if err != nil {
		fmt.Fprintln(w, "sign", err.Error())
		return
	}
	fmt.Fprintf(w, `{"token": "%s"}\n`, tokenStr)
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
	channel := content.Dishes.Channel
	fmt.Println(channel, data["content"])
	rep, err := http.NewRequest("POST", API_SEND+channel+"/messages", bytes.NewBuffer(dataBytes))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	rep.Header.Set("Authorization", "Bot "+token)
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
		"postgresql://%s:%s@database:5432/area_database?sslmode=disable",
		dbUser,
		dbPassword,
	)
	return sql.Open("postgres", connectStr)
}

func miniproxy(f func(http.ResponseWriter, *http.Request, *sql.DB), c *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(a http.ResponseWriter, b *http.Request) {
		f(a, b, c)
	}
}

type Spice struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Route struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Spices []Spice `json:"spices"`
}

func getRoutes(w http.ResponseWriter, req *http.Request) {
	var list = []Route{
		Route{
			Name: "send",
			Type: "reaction",
			Spices: []Spice{
				{
					Name: "channel",
					Type: "number",
				},
				{
					Name: "message",
					Type: "text",
				},
			},
		},
	}
	var data []byte
	var err error

	data, err = json.Marshal(list)
	if err != nil {
		http.Error(w, `{ "error":  "marshal" }`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, string(data))
}

func main() {
	db, err := connectToDatabase()
	if err != nil {
		os.Exit(84)
	}
	fmt.Println("discord microservice container is running !")
	router := mux.NewRouter()
	godotenv.Load("/usr/mound.d/.env", "/usr/mound.d/.env1")
	router.HandleFunc("/send", doSomeSend).Methods("POST")
	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	router.HandleFunc("/routes", getRoutes).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}
