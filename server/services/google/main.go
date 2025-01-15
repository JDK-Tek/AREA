package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
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

const API_OAUTH_GOOGLE = "https://accounts.google.com/o/oauth2/token"
const API_USER_GOOGLE = "https://www.googleapis.com/oauth2/v2/userinfo"
const API_SEND_GOOGLE = "https://www.googleapis.com/upload/gmail/v1/users/me/messages/send"

const PERMISSIONS = 8
const EXPIRATION = 60 * 30

type Result struct {
	Code string `json:"code"`
}

type TokenResult struct {
	Token   string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

type UserResult struct {
	ID string `json:"id"`
}

type EmailContent struct {
	UserID  string `json:"UserID"`
	Token   string `json:"Token"`
	Subject string `json:"Subject"`
	Body    string `json:"Body"`
	To      string `json:"To"`
}

func getOAUTHLink(w http.ResponseWriter, req *http.Request) {
	params := url.Values{}
	params.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	params.Set("response_type", "code")
	params.Set("redirect_uri", os.Getenv("REDIRECT"))
	params.Set("scope", "openid profile email offline_access Gmail.Send")
	params.Set("state", "some-state-value")

	oauthURL := "https://accounts.google.com/o/oauth2/auth?" + params.Encode()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, oauthURL)
}

func setOAUTHToken(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var res Result
	var tok TokenResult
	var user UserResult
	var responseData map[string]interface{}

	clientid := os.Getenv("GOOGLE_CLIENT_ID")
	clientsecret := os.Getenv("GOOGLE_CLIENT_SECRET")
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
	rep, err := http.PostForm(API_OAUTH_GOOGLE, data)
	body, err := io.ReadAll(rep.Body)
	if err != nil {
		fmt.Fprintln(w, "postform", err.Error())
		return
	}
	defer rep.Body.Close()

	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Fprintln(w, "unmarshal json", err.Error())
		return
	}

	tok.Token = responseData["access_token"].(string)
	tok.Refresh = responseData["refresh_token"].(string)

	req, err = http.NewRequest("GET", API_USER_GOOGLE, nil)
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

	if tok.Token == "" || tok.Refresh == "" {
		fmt.Fprintln(w, "error: token is empty")
		return
	}

	var tokid int
	err = db.QueryRow("select id from tokens where userid = $1", user.ID).Scan(&tokid)
	if err != nil {
		err = db.QueryRow("insert into tokens (service, token, userid) values ($1, $2, $3) returning id",
			"google", tok.Token, user.ID).Scan(&tokid)
		if err != nil {
			fmt.Fprintln(w, "db insert", err.Error())
			return
		}
	}

	secretBytes := []byte(os.Getenv("BACKEND_KEY"))
	claims := jwt.MapClaims{
		"id":  user.ID,
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
	var emailContent EmailContent
	var tok TokenResult

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&emailContent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	userID := emailContent.UserID
	tokenStr := emailContent.Token

	if tokenStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{ \"error\": \"Token is missing\" }\n")
		return
	}

	secretBytes := []byte(os.Getenv("BACKEND_KEY"))
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretBytes, nil
	})

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{ \"error\": \"Invalid token\" }\n")
		return
	}

	err = db.QueryRow("SELECT token, refresh FROM tokens WHERE userid = $1", userID).Scan(&tok.Token, &tok.Refresh)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "{ \"error\": \"Token not found for the user\" }\n")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{ \"error\": \"Database query error: %s\" }\n", err.Error())
		}
		return
	}

	emailData := map[string]interface{}{
		"raw": encodeWeb64(emailContent.Subject, emailContent.Body, emailContent.To),
	}

	emailBytes, err := json.Marshal(emailData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"Error marshalling email data: %s\" }\n", err.Error())
		return
	}

	reqEmail, err := http.NewRequest("POST", API_SEND_GOOGLE, bytes.NewBuffer(emailBytes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"Error creating request: %s\" }\n", err.Error())
		return
	}
	reqEmail.Header.Set("Authorization", "Bearer " + tok.Token)
	client := &http.Client{}
	rep, err := client.Do(reqEmail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"Error sending email: %s\" }\n", err.Error())
		return
	}
	defer rep.Body.Close()

	if rep.StatusCode == http.StatusOK {
		fmt.Fprintf(w, `{ "message": "Email sent successfully" }\n`)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"Failed to send email\" }\n")
	}
}

func encodeWeb64(subject, body, to string) string {
	message := fmt.Sprintf("Subject: %s\nTo: %s\n\n%s", subject, to, body)
	encoded := base64.URLEncoding.EncodeToString([]byte(message))
	return encoded
}

func main() {
	godotenv.Load()

	db, err := sql.Open("postgres", os.Getenv("DB_CONNECTION"))
	if err != nil {
		log.Fatal("Erreur de connexion à la base de données", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/oauth/google", getOAUTHLink).Methods("GET")
	r.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		setOAUTHToken(w, r, db)
	}).Methods("POST")
	r.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		sendEmail(w, r, db)
	}).Methods("POST")

	log.Fatal(http.ListenAndServe(":80", r))
}
