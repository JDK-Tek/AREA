package main

import (
	//"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
	"strings"
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
	fmt.Println("test test")
    str := "https://accounts.google.com/o/oauth2/v2/auth?"
    
    redirectURI := url.QueryEscape(os.Getenv("REDIRECT"))
    fmt.Println("Redirect URI = ", redirectURI)

    scopes := "https://www.googleapis.com/auth/drive.file " +
              "https://www.googleapis.com/auth/userinfo.profile " +
              "https://www.googleapis.com/auth/userinfo.email " +
              "https://www.googleapis.com/auth/gmail.send"

    str += "client_id=" + os.Getenv("GOOGLE_CLIENT_ID")
    str += "&response_type=code"
    str += "&redirect_uri=" + redirectURI
    str += "&scope=" + url.QueryEscape(scopes)
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
	fmt.Println("i m here")
    var res Result
    var tok TokenResult
    var user UserResult
    var tokid int
    var owner = -1
    var responseData map[string]interface{}

    clientID := os.Getenv("GOOGLE_CLIENT_ID")
    clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
    redirectURI := os.Getenv("REDIRECT")
    data := url.Values{}
    
    err := json.NewDecoder(req.Body).Decode(&res)
    if err != nil {
        fmt.Fprintln(w, "Erreur lors du décodage de la requête:", err.Error())
        return
    }

    data.Set("client_id", strings.TrimSpace(clientID))
    data.Set("client_secret", strings.TrimSpace(clientSecret))
    data.Set("grant_type", "authorization_code")
    data.Set("code", strings.TrimSpace(url.QueryEscape(res.Code)))
    data.Set("redirect_uri", strings.TrimSpace(redirectURI))
    
	fmt.Println("code = ", res.Code)
	fmt.Println("client secret = ", clientSecret)
	fmt.Println("so redirect = ", redirectURI)

    rep, err := http.PostForm("https://oauth2.googleapis.com/token", data)
    if err != nil {
        fmt.Fprintln(w, "Erreur lors de l'échange du code:", err.Error())
        return
    }
    defer rep.Body.Close()
    
    body, err := io.ReadAll(rep.Body)
	fmt.Println("rep = ", string(body))
	fmt.Println("status = ", rep.StatusCode)
    if err != nil {
        fmt.Fprintln(w, "Erreur lors de la lecture du corps de la réponse:", err.Error())
		fmt.Println("error: ", err.Error())
        return
    }
    
    if err := json.Unmarshal(body, &responseData); err != nil {
        fmt.Fprintln(w, "Erreur lors de l'analyse de la réponse JSON:", err.Error())
        return
    }

    if tokenStr, ok := responseData["access_token"].(string); ok {
        tok.Token = tokenStr
    } else {
        fmt.Fprintln(w, "Erreur : access_token manquant ou incorrect")
        return
    }

    if refreshTokenStr, ok := responseData["refresh_token"].(string); ok {
        tok.Refresh = refreshTokenStr
    } else {
        fmt.Fprintln(w, "Erreur : refresh_token manquant ou incorrect")
        return
    }

    req, err = http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
    if err != nil {
        fmt.Fprintln(w, "Erreur lors de la création de la requête utilisateur:", err.Error())
        return
    }
    req.Header.Set("Authorization", "Bearer "+tok.Token)
    
    client := &http.Client{}
    rep, err = client.Do(req)
    if err != nil {
        fmt.Fprintln(w, "Erreur lors de l'appel à l'API Google:", err.Error())
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
            "google", tok.Token, tok.Refresh, user.ID).Scan(&tokid)
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
    
    fmt.Println("Succès de l'authentification avec Google, token =", tokenStr)
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
		http.Error(w, "Erreur lors de l'appel API Google", http.StatusBadGateway)
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

type Message struct {
    Bridge int `json:"bridge"`
    UserId int `json:"userid"`
    Ingredients map[string]string `json:"ingredients"`
}


func getRoutes(w http.ResponseWriter, req *http.Request) {
	var list = []InfoRoute{
		{
			Name: "playMusic",
			Type: "reaction",
			Desc: "Chose a music to play !",
			Spices: []InfoSpice{
				{
					Name: "musique",
					Type: "text",
					Title: "The message you want to send.",
				},
			},
		},
		{
			Name: "pauseMusic",
			Type: "reaction",
			Desc: "Stop the current music !",
			Spices: []InfoSpice{},
		},
		{
			Name: "resumeMusic",
			Type: "reaction",
			Desc: "Resume the current music !",
			Spices: []InfoSpice{},
		},
		{
			Name: "checkDeviceConnection",
			Type: "action",
			Desc: "check if you have a current spotify running !",
			Spices: []InfoSpice{},
		},
		{
			Name: "checkSongRunning",
			Type: "action",
			Desc: "check if you have a song whose running and got id, type and current song !",
			Spices: []InfoSpice{},
		},
		{
			Name: "checkPodcastRunning",
			Type: "action",
			Desc: "check if you have a podcast whose running and got the podcast name and episode + follow him !",
			Spices: []InfoSpice{},
		},
		{
			Name: "likeCurrentMusic",
			Type: "reaction",
			Desc: "like the current music you listening",
			Spices: []InfoSpice{},
		},
		{
			Name: "unlikeCurrentMusic",
			Type: "reaction",
			Desc: "unlike the current music you listening",
			Spices: []InfoSpice{},
		},
		{
			Name: "nextMusic",
			Type: "reaction",
			Desc: "next music",
			Spices: []InfoSpice{},
		},
		{
			Name: "previousMusic",
			Type: "reaction",
			Desc: "previous music",
			Spices: []InfoSpice{},
		},
		{
			Name: "removeFromPlaylistIfPresent",
			Type: "reaction",
			Desc: "remove the current music if she are in a playlist",
			Spices: []InfoSpice{},
		},
		{
			Name: "addToPlaylistIfNotPresent",
			Type: "reaction",
			Desc: "add the current music if she are in a playlist",
			Spices: []InfoSpice{},
		},
		{
			Name: "createPlaylist",
			Type: "reaction",
			Desc: "creer une playlist",
			Spices: []InfoSpice{
				{
					Name: "name",
					Type: "text",
					Title: "The title for the playlist",
				},
			},
		},
		{
			Name: "clearPlaylist",
			Type: "reaction",
			Desc: "clear une playlist",
			Spices: []InfoSpice{
				{
					Name: "name",
					Type: "text",
					Title: "The title of the playlist",
				},
			},
		},
		{
			Name: "addToPlaylistByName",
			Type: "reaction",
			Desc: "add the current song listening to the playlist selected",
			Spices: []InfoSpice{
				{
					Name: "name",
					Type: "text",
					Title: "The title of the playlist",
				},
			},
		},
	}
	var infos = Infos{
		Color: "#1DB954",
		Image: "https://m.media-amazon.com/images/I/51rttY7a+9L.png",
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
	fmt.Println("Google microservice container is running !")
	router := mux.NewRouter()
	godotenv.Load(".env")

	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	router.HandleFunc("/user", getUserInfo).Methods("GET")
	router.HandleFunc("/", getRoutes).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}