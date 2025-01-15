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

	"golang.org/x/oauth2"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"encoding/base64"
)

const API_OAUTH_CUSTOM = "https://example.com/oauth/token" // Remplacer par l'URL de l'API appropriée
const API_USER_CUSTOM = "https://example.com/api/v1/user" // Remplacer par l'URL de l'API appropriée
const API_EMAIL_SEND = "https://example.com/api/v1/send-email" // Remplacer par l'URL de l'API appropriée

const EXPIRATION = 60 * 30

type Result struct {
	Code string `json:"code"`
}

type TokenResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

type UserResult struct {
	ID string `json:"sub"`
}

func getOAUTHLink(w http.ResponseWriter, req *http.Request) {
	str := "https://example.com/oauth/authorize?" // Remplacer par l'URL de l'API appropriée
	str += "client_id=" + os.Getenv("CUSTOM_CLIENT_ID")
	str += "&response_type=code"
	str += "&redirect_uri=" + url.QueryEscape(os.Getenv("REDIRECT"))
	str += "&scope=openid profile email send_email"
	str += "&state=some-state-value"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, str)
}

func setOAUTHToken(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var res Result
	var tok TokenResult
	var user UserResult
	var tokid int
	var owner = -1

	clientid := os.Getenv("CUSTOM_CLIENT_ID")
	clientsecret := os.Getenv("CUSTOM_CLIENT_SECRET")
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

	resp, err := http.PostForm(API_OAUTH_CUSTOM, data)
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

	if tok.AccessToken == "" || tok.RefreshToken == "" {
		fmt.Fprintln(w, "error: token is empty")
		return
	}

	req, err = http.NewRequest("GET", API_USER_CUSTOM, nil)
	if err != nil {
		fmt.Fprintln(w, "request error", err.Error())
		return
	}
	req.Header.Set("Authorization", "Bearer " + tok.AccessToken)
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
			"custom",
			tok.AccessToken,
			tok.RefreshToken,
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

func sendEmail(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var res struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	err := json.NewDecoder(req.Body).Decode(&res)
	if err != nil {
		fmt.Fprintln(w, "decode", err.Error())
		return
	}

	token := req.Header.Get("Authorization")
	if token == "" {
		fmt.Fprintln(w, "Authorization header is missing")
		return
	}

	tok := &oauth2.Token{AccessToken: token}

	service, err := gmail.NewService(context.Background(), option.WithTokenSource(oauth2.StaticTokenSource(tok)))
	if err != nil {
		fmt.Fprintln(w, "failed to create custom email service", err.Error())
		return
	}

	message := &gmail.Message{}

	emailTo := []string{res.To}
	emailSubject := res.Subject
	emailBody := res.Body

	emailMessage := []byte(fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", "user@example.com", emailTo[0], emailSubject, emailBody))

	encodedMessage := base64.URLEncoding.EncodeToString(emailMessage)

	message.Raw = encodedMessage

	_, err = service.Users.Messages.Send("me", message).Do()
	if err != nil {
		fmt.Fprintln(w, "Failed to send email", err.Error())
		return
	}

	fmt.Fprintln(w, "Email sent successfully")
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
	fmt.Println("Custom microservice container is running !")
	router := mux.NewRouter()
	godotenv.Load(".env")
	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	router.HandleFunc("/send-email", miniproxy(sendEmail, db)).Methods("POST")
	log.Fatal(http.ListenAndServe(":80", router))
}
