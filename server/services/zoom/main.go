package main

import (
	//"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	//"io"
	"log"
	"net/http"
	"net/url"
	"os"
	//"time"

	//"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	//_ "github.com/lib/pq"
)

const (
	ZOOM_AUTH_URL   = "https://zoom.us/oauth/authorize"
	ZOOM_TOKEN_URL  = "https://zoom.us/oauth/token"
	ZOOM_API_ME_URL = "https://api.zoom.us/v2/users/me"
)

const PERMISSIONS = 8
const EXPIRATION = 60 * 30

var db *sql.DB

type InfoSpice struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	Title string   `json:"title"`
	Extra []string `json:"extra"`
}

type InfoRoute struct {
	Type   string      `json:"type"`
	Name   string      `json:"name"`
	Desc   string      `json:"description"`
	Spices []InfoSpice `json:"spices"`
}

type Infos struct {
	Color  string      `json:"color"`
	Image  string      `json:"image"`
	Routes []InfoRoute `json:"areas"`
}

func getUserInfo(w http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Token d'accès manquant", http.StatusUnauthorized)
		return
	}

	reqUser, err := http.NewRequest("GET", ZOOM_API_ME_URL, nil)
	if err != nil {
		http.Error(w, "Erreur lors de la création de la requête utilisateur", http.StatusInternalServerError)
		return
	}

	reqUser.Header.Set("Authorization", "Bearer "+authHeader[7:])

	client := &http.Client{}
	resp, err := client.Do(reqUser)
	if err != nil {
		http.Error(w, "Erreur lors de l'appel API Spotify", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Erreur lors de la décodification des informations utilisateur", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userInfo)
}

func getOAUTHLink(w http.ResponseWriter, req *http.Request) {
    fmt.Println("test test")
    str := "https://zoom.us/oauth/authorize?"

    redirectURI := url.QueryEscape(os.Getenv("REDIRECT"))
    fmt.Println("Redirect URI = ", redirectURI)

    scopes := "user:read " +
              "meeting:write " +
              "user:write " +
              "account:read"

    str += "client_id=" + os.Getenv("ZOOM_CLIENT_ID")
    str += "&response_type=code"
    str += "&scope=" + url.QueryEscape(scopes)
    str += "&state=some-state-value"
    str += "&redirect_uri=" + redirectURI

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, str)
}

func getRoutes(w http.ResponseWriter, req *http.Request) {
	var list = []InfoRoute{
		{
			Name: "playMusic",
			Type: "reaction",
			Desc: "Chose a music to play !",
			Spices: []InfoSpice{
				{
					Name:  "musique",
					Type:  "text",
					Title: "The message you want to send.",
				},
			},
		},
		{
			Name: "pauseMusic",
			Type: "reaction",
			Desc: "Stop the current music !",
			Spices: []InfoSpice{
				{
					Name:  "message",
					Type:  "text",
					Title: "The message you want to send.",
				},
			},
		},
	}
	var infos = Infos{
		Color:  "#1DB954",
		Image:  "/assets/spotify.webp",
		Routes: list,
	}
	var data []byte
	var err error

	data, err = json.Marshal(infos)
	if err != nil {
		http.Error(w, `{ "error":  "marshal" }`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(data))
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

func main() {
	// db, err := connectToDatabase()
	// if err != nil {
	// 	os.Exit(84)
	// }
	fmt.Println("Zoom microservice container is running !")
	router := mux.NewRouter()
	godotenv.Load(".env")

	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	//router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	router.HandleFunc("/user", getUserInfo).Methods("GET")
	router.HandleFunc("/", getRoutes).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}