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
    str := "https://accounts.google.com/o/oauth2/v2/auth?"
    
    redirectURI := url.QueryEscape(os.Getenv("REDIRECT_URI"))
	fmt.Println("test test")
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

func setOAUTHToken(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var res Result
	var tok TokenResult
	var user UserResult
	var tokid int
	var owner = -1
	var responseData map[string]interface{}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT_URI")
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
	
	rep, err := http.PostForm("https://oauth2.googleapis.com/token", data)
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
			"google",
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
	
	fmt.Println("Succès de l'authentification avec Google, token =", tokenStr)
	fmt.Fprintf(w, "{\"token\": \"%s\"}\n", tokenStr)
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
	fmt.Println("Google microservice container is running !")
	router := mux.NewRouter()
	godotenv.Load(".env")

	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	log.Fatal(http.ListenAndServe(":80", router))
}