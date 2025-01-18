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
    fmt.Println("redirect = ", x)
    
    str += "client_id=" + os.Getenv("SPOTIFY_CLIENT_ID")
    str += "&response_type=code"
    str += "&redirect_uri=" + x
    str += "&scope=" +
        "user-library-read " +
        "playlist-read-private " +
        "playlist-read-collaborative " +
        "user-read-playback-state " +
        "user-read-currently-playing " +
        "user-modify-playback-state " +
        "app-remote-control " +
        "user-top-read " +
        "playlist-modify-public " +
        "playlist-modify-private " +
        "streaming " +
        "user-follow-read " +
        "user-follow-modify"
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

func playMusic(w http.ResponseWriter, req *http.Request, db *sql.DB) {
    fmt.Println("Headers received:", req.Header)

    bodyBytes, err := io.ReadAll(req.Body)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error reading request body:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }
    fmt.Println("Request Body:", string(bodyBytes))

    req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

    var requestBody struct {
        UserID int `json:"userid"`
    }

    decoder := json.NewDecoder(req.Body)
    err = decoder.Decode(&requestBody)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error decoding JSON:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }

    userID := requestBody.UserID
    userID = 1
    fmt.Println("Extracted userID:", userID)

    var spotifyToken string
    err = db.QueryRow("SELECT token FROM tokens WHERE owner = $1 AND service = 'spotify'", userID).Scan(&spotifyToken)
    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Println("No Spotify token found for user:", userID)
            fmt.Fprintf(w, "{ \"error\": \"No Spotify token found for user\" }\n")
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Println("Database error:", err.Error())
            fmt.Fprintf(w, "{ \"error\": \"Database error: %s\" }\n", err.Error())
        }
        return
    }

    if spotifyToken == "" {
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Println("No Spotify token available for user:", userID)
        fmt.Fprintf(w, "{ \"error\": \"No Spotify token available\" }\n")
        return
    }

    trackURI := "spotify:track:3n3P1vEXs6IfzozT8kVYAf"
    spotifyURL := "https://api.spotify.com/v1/me/player/play"
    body := fmt.Sprintf(`{"uris":["%s"]}`, trackURI)

    reqSpotify, err := http.NewRequest("PUT", spotifyURL, bytes.NewBufferString(body))
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error creating Spotify request:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }

    reqSpotify.Header.Set("Authorization", "Bearer "+spotifyToken)
    reqSpotify.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    respSpotify, err := client.Do(reqSpotify)
    if err != nil {
        w.WriteHeader(http.StatusBadGateway)
        fmt.Println("Error playing music on Spotify:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }
    defer respSpotify.Body.Close()

    if respSpotify.StatusCode == http.StatusNoContent {
        fmt.Println("Music 'We Are the Champions' is now playing!")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "{ \"status\": \"Music is now playing!\" }\n")
    } else {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Failed to play music on Spotify. Status:", respSpotify.StatusCode)
		fmt.Println("rep = ", reqSpotify.Body)
        fmt.Fprintf(w, "{ \"error\": \"Failed to play music\" }\n")
    }
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

type InfoSpice struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Title string `json:"title"`
	Extra []string `json:"extra"`
}

type InfoRoute struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Desc string `json:"description"`
	Spices []InfoSpice `json:"spices"`
}

type Infos struct {
	Color string `json:"color"`
	Image string `json:"image"`
	Routes []InfoRoute `json:"areas"`
}

func getRoutes(w http.ResponseWriter, req *http.Request) {
	var list = []InfoRoute{
		InfoRoute{
			Name: "playMusic",
			Type: "reaction",
			Desc: "Sends an email.",
			Spices: []InfoSpice{
				{
					Name: "musique",
					Type: "text",
					Title: "The message you want to send.",
				},
			},
		},
	}
	var infos = Infos{
		Color: "#5865F2",
		Image: "https://cdn.prod.website-files.com/6257adef93867e50d84d30e2/636e0a6cc3c481a15a141738_icon_clyde_white_RGB.png",
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

func main() {
	db, err := connectToDatabase()
	if err != nil {
		os.Exit(84)
	}
	fmt.Println("Spotify microservice container is running !")
	router := mux.NewRouter()
	godotenv.Load(".env")

	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	router.HandleFunc("/playMusic", miniproxy(playMusic, db)).Methods("POST")
	router.HandleFunc("/user", getUserInfo).Methods("GET")
	router.HandleFunc("/", getRoutes).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}
