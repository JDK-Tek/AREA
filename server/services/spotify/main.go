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

const API_OAUTH_OUTLOOK = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
const API_USER_OUTLOOK = "https://graph.microsoft.com/v1.0/me"
const API_SEND = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
const API_OAUTH = "https://discord.com/api/oauth2/token"
const API_USER = "https://discord.com/api/v10/users/@me"

const PERMISSIONS = 8
const EXPIRATION = 60 * 30

// Structures pour envoyer des messages Teams ou e-mails
type TeamsContent struct {
	TeamID   string `json:"team_id"`
	ChannelID string `json:"channel_id"`
	Message  string `json:"message"`
}

type EmailContent struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func getOAUTHLink(w http.ResponseWriter, req *http.Request) {
	params := url.Values{}
	params.Set("client_id", os.Getenv("OUTLOOK_CLIENT_ID"))
	params.Set("response_type", "code")
	params.Set("redirect_uri", os.Getenv("REDIRECT"))
	params.Set("scope", "openid profile email offline_access user.read")
	params.Set("state", "some-state-value")

	oauthURL := "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?" + params.Encode()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, oauthURL)
}

type Result struct {
	Code string `json:"code"`
}

type TokenResult struct {
	Token  string `json:"access_token"`
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

	clientid := os.Getenv("OUTLOOK_CLIENT_ID")
	clientsecret := os.Getenv("OUTLOOK_CLIENT_SECRET")
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
	rep, err := http.PostForm(API_OAUTH_OUTLOOK, data)
	body, err := io.ReadAll(rep.Body)
	fmt.Println(clientsecret)
	fmt.Fprintln(w, "tmp = ", string(body))
	if err != nil {
		fmt.Fprintln(w, "postform", err.Error())
		return
	}
	defer rep.Body.Close()

	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Fprintln(w, "unmarshal json", err.Error())
		return
	}
	fmt.Println(responseData)
	tok.Token = responseData["access_token"].(string)
	tok.Refresh = "foobar"

	req, err = http.NewRequest("GET", API_USER_OUTLOOK, nil)
	if err != nil {
		fmt.Fprintln(w, "request error", err.Error())
		return
	}
	req.Header.Set("Authorization", "Bearer " + tok.Token)
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

	// Token handling for database storage
	if tok.Token == "" || tok.Refresh == "" {
		fmt.Fprintln(w, "error: token is empty")
		return
	}
	err = db.QueryRow("select id, owner from tokens where userid = $1", user.ID).Scan(&tokid, &owner)
	if err != nil {
		err = db.QueryRow("insert into tokens (service, token, userid) values ($1, $2, $3) returning id",
			"outlook",
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

// Send email via Outlook API
func sendEmail(w http.ResponseWriter, req *http.Request) {
	var emailContent EmailContent
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&emailContent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	token := os.Getenv("OUTLOOK_TOKEN")
	if token == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"token is missing\" }\n")
		return
	}

	emailData := map[string]interface{}{
		"message": map[string]interface{}{
			"subject": emailContent.Subject,
			"body": map[string]interface{}{
				"contentType": "Text",
				"content":     emailContent.Body,
			},
			"toRecipients": []map[string]interface{}{
				{
					"emailAddress": map[string]string{
						"address": emailContent.To,
					},
				},
			},
		},
		"saveToSentItems": "true",
	}
	emailBytes, err := json.Marshal(emailData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqEmail, err := http.NewRequest("POST", "https://graph.microsoft.com/v1.0/me/sendMail", bytes.NewBuffer(emailBytes))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqEmail.Header.Set("Authorization", "Bearer "+token)
	reqEmail.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqEmail)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{ \"status\": \"Email sent successfully\" }\n")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"Failed to send email\" }\n")
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

func miniproxy(f func(http.ResponseWriter, *http.Request, *sql.DB), c *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(a http.ResponseWriter, b *http.Request) {
		f(a, b, c)
	}
}

func main() {
	godotenv.Load(".env")
	db, err := connectToDatabase()
	if err != nil {
		os.Exit(84)
	}
	fmt.Println("Outlook microservice is running!")
	router := mux.NewRouter()
	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	router.HandleFunc("/sendEmail", sendEmail).Methods("POST")
	log.Fatal(http.ListenAndServe(":80", router))
}
