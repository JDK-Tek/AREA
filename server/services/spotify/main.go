package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
	"io"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)
const (
	API_OAUTH_SPOTIFY = "https://accounts.spotify.com/api/token"
	API_USER_SPOTIFY  = "https://api.spotify.com/v1/me"
)

const PERMISSIONS = 8
const EXPIRATION = 60 * 30


var db *sql.DB

func init() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
}

func getOAUTHLink(w http.ResponseWriter, req *http.Request) {
	str := "https://accounts.spotify.com/authorize?"
	x := url.QueryEscape(os.Getenv("REDIRECT"))
	str += "client_id=" + os.Getenv("SPOTIFY_CLIENT_ID")
	str += "&response_type=code"
	str += "&redirect_uri=" + x
	str += "&scope=user-library-read"
	str += "&state=some-state-value"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, str)
}

type Result struct {
	Code string `json:"code"`
}

type TokenResult struct {
	Token string `json:"access_token"`
	Refresh string `json:"refresh_token"`
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
	var responseData map[string]interface{}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT")
	data := url.Values{}
	
	err := json.NewDecoder(req.Body).Decode(&res)
	if err != nil {
		fmt.Fprintln(w, "Erreur lors du décodage de la requête:", err.Error())
		return
	}

	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", res.Code)
	data.Set("redirect_uri", redirectURI)
	
	rep, err := http.PostForm("https://accounts.spotify.com/api/token", data)
	if err != nil {
		fmt.Fprintln(w, "Erreur lors de l'échange du code:", err.Error())
		return
	}
	defer rep.Body.Close()
	
	body, err := io.ReadAll(rep.Body)
	if err != nil {
		fmt.Fprintln(w, "Erreur lors de la lecture du corps de la réponse:", err.Error())
		return
	}
	
	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Fprintln(w, "Erreur lors de l'analyse de la réponse JSON:", err.Error())
		return
	}

	tok.Token = responseData["access_token"].(string)
	tok.Refresh = responseData["refresh_token"].(string)

	req, err = http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		fmt.Fprintln(w, "Erreur lors de la création de la requête utilisateur:", err.Error())
		return
	}
	req.Header.Set("Authorization", "Bearer "+tok.Token)
	
	client := &http.Client{}
	rep, err = client.Do(req)
	if err != nil {
		fmt.Fprintln(w, "Erreur lors de l'appel à l'API Spotify:", err.Error())
		return
	}
	defer rep.Body.Close()
	
	err = json.NewDecoder(rep.Body).Decode(&user)
	if err != nil {
		fmt.Fprintln(w, "Erreur lors du décodage des informations utilisateur:", err.Error())
		return
	}

	if tok.Token == "" || tok.Refresh == "" {
		fmt.Fprintln(w, "Erreur : token ou refresh token manquant")
		return
	}

	err = db.QueryRow("SELECT id, owner FROM tokens WHERE userid = $1", user.ID).Scan(&tokid, &owner)
	if err != nil {
		err = db.QueryRow("INSERT INTO tokens (service, token, refresh, userid) VALUES ($1, $2, $3, $4) RETURNING id",
			"spotify",
			tok.Token,
			tok.Refresh,
			user.ID,
		).Scan(&tokid)
		if err != nil {
			fmt.Fprintln(w, "Erreur lors de l'insertion du token:", err.Error())
			return
		}

		err = db.QueryRow("INSERT INTO users (tokenid) VALUES ($1) RETURNING id", tokid).Scan(&owner)
		if err != nil {
			fmt.Fprintln(w, "Erreur lors de l'insertion de l'utilisateur:", err.Error())
			return
		}
		db.Exec("UPDATE tokens SET owner = $1 WHERE id = $2", owner, tokid)
	}

	secretBytes := []byte(os.Getenv("BACKEND_KEY"))
	claims := jwt.MapClaims{
		"id":  owner,
		"exp": time.Now().Add(time.Second * EXPIRATION).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secretBytes)
	if err != nil {
		fmt.Fprintln(w, "Erreur lors de la signature du token:", err.Error())
		return
	}
	
	fmt.Println("Succès de l'authentification avec Spotify, token =", tokenStr)
	fmt.Fprintf(w, "{\"token\": \"%s\"}\n", tokenStr)
}

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

func miniproxy(f func(http.ResponseWriter, *http.Request, *sql.DB), c *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(a http.ResponseWriter, b *http.Request) {
		f(a, b, c)
	}
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()

	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	router.HandleFunc("/user", getUserInfo).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
