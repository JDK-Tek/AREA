package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	ZOOM_AUTH_URL   = "https://zoom.us/oauth/authorize"
	ZOOM_TOKEN_URL  = "https://zoom.us/oauth/token"
	ZOOM_API_ME_URL = "https://api.zoom.us/v2/users/me"
)

const PERMISSIONS = 8
const EXPIRATION = 60 * 30

func getUserInfo(w http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Token d'accès manquant", http.StatusUnauthorized)
		return
	}

	reqUser, err := http.NewRequest("GET", API_USER_SPOTIFY, nil)
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
    str := "https://zoom.us/oauth/authorize"

    redirectURI := url.QueryEscape(os.Getenv("REDIRECT"))
    fmt.Println("Redirect URI = ", redirectURI)

    scopes := "user:read " +
              "meeting:write " +
              "user:write " +
              "account:read"

    str += "client_id=" + os.Getenv("ZOOM_CLIENT_ID")
    str += "&response_type=code"
    str += "&redirect_uri=" + redirectURI
    str += "&scope=" + url.QueryEscape(scopes)
    str += "&state=some-state-value"

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, str)
}

func main() {
	db, err := connectToDatabase()
	if err != nil {
		os.Exit(84)
	}
	fmt.Println("Zoom microservice container is running !")
	router := mux.NewRouter()
	godotenv.Load(".env")

	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	//router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	router.HandleFunc("/user", getUserInfo).Methods("GET")
	router.HandleFunc("/", getRoutes).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}