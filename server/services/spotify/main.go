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

	// "io/ioutil"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const API_OAUTH_SPOTIFY = "https://accounts.spotify.com/api/token"
const API_USER_SPOTIFY = "https://api.spotify.com/v1/me"

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
	str := "https://accounts.spotify.com/authorize?"
	str += "client_id=" + os.Getenv("SPOTIFY_CLIENT_ID")
	str += "&response_type=code"
	str += "&redirect_uri=" + url.QueryEscape(os.Getenv("REDIRECT"))
	str += "&scope=user-read-private user-read-email"
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

    clientid := os.Getenv("SPOTIFY_CLIENT_ID")
    clientsecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
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
    data.Set("redirect_uri", os.Getenv("REDIRECT"))

    resp, err := http.PostForm(API_OAUTH_SPOTIFY, data)
    if err != nil {
        fmt.Fprintln(w, "postform: ", err.Error())
        return
    }
    defer resp.Body.Close()

    err = json.NewDecoder(resp.Body).Decode(&tok)
    if err != nil {
        fmt.Fprintln(w, "decode: ", err.Error())
        return
    }

    if tok.Token == "" || tok.Refresh == "" {
        fmt.Fprintln(w, "error: token is empty")
        return
    }
    
    req, err = http.NewRequest("GET", API_USER_SPOTIFY, nil)
    if err != nil {
        fmt.Fprintln(w, "request error", err.Error())
        return
    }
    req.Header.Set("Authorization", "Bearer " + tok.Token)
    client := &http.Client{}
    resp, err = client.Do(req)
    if err != nil {
        fmt.Fprintln(w, "client do", err.Error())
        return
    }
    defer resp.Body.Close()

    err = json.NewDecoder(resp.Body).Decode(&user)
    if err != nil {
        fmt.Fprintln(w, "decode", err.Error())
        return
    }

    err = db.QueryRow("select id, owner from tokens where userid = $1", user.ID).Scan(&tokid, &owner)
    if err != nil {
        err = db.QueryRow("insert into tokens (service, token, refresh, userid) values ($1, $2, $3, $4) returning id",
            "spotify",
            tok.Token,
            tok.Refresh,
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

func miniproxy(f func(http.ResponseWriter, *http.Request, *sql.DB), c *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(a http.ResponseWriter, b *http.Request) {
		f(a, b, c)
	}
}

func main() {
	db, err := connectToDatabase()
	if err != nil {
		os.Exit(84)
	}
	fmt.Println("spotify microservice container is running !")
	router := mux.NewRouter()
	godotenv.Load(".env")
	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	log.Fatal(http.ListenAndServe(":80", router))
}
