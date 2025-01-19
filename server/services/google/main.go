package main

import (
	"bytes"
	"context"
	"database/sql"
	"strconv"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

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
	Token   string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

type UserResult struct {
	ID string `json:"id"`
}

func getIdFromToken(tokenString string) (int, error) {
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = tokenString[len("Bearer "):]
	}
	fmt.Println(tokenString)
	secretKey := []byte(os.Getenv("BACKEND_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return -1, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["id"].(string)
		if !ok {
			return -1, fmt.Errorf("'id' field not found or not a string")
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return -1, fmt.Errorf("error converting id to int: %v", err)
		}
		return idInt, nil
	} else {
		return -1, fmt.Errorf("invalid token")
	}
}

func setOAUTHToken(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	fmt.Println("i m here")
	var res Result
	var tok TokenResult
	var user UserResult
	var tokid int
	var owner = -1
	var responseData map[string]interface{}
	var atok = req.Header.Get("Authorization")

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT")
	data := url.Values{}

	err := json.NewDecoder(req.Body).Decode(&res)
	if err != nil {
		fmt.Fprintln(w, "Erreur lors du décodage de la requête:", err.Error())
		return
	}

	decodedCode, err := url.QueryUnescape(res.Code)
	if err != nil {
		return
	}
	data.Set("client_id", strings.TrimSpace(clientID))
	data.Set("client_secret", strings.TrimSpace(clientSecret))
	data.Set("grant_type", "authorization_code")
	data.Set("code", decodedCode)
	data.Set("redirect_uri", strings.TrimSpace(redirectURI))

	fmt.Println("code = ", decodedCode)
	fmt.Println("client secret = ", strings.TrimSpace(clientSecret))
	fmt.Println("so redirect = ", strings.TrimSpace(redirectURI))

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

	tok.Token = responseData["access_token"].(string)

	fmt.Println("acess token = ", tok.Token)
	fmt.Println("refresk token = ", "")

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

	if tok.Token == "" {
		fmt.Fprintln(w, "Erreur : token ou refresh token manquant")
		return
	}

	// inserting into database, first i get the 'users' token
	if atok != "" {
		// if the user is logged, i get the userid
		tokid, err = getIdFromToken(tok.Token)
		if err != nil {
			fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
			return
		}
	} else {
		// if the user is not logged, create an empty user
		err = db.QueryRow("insert into users default values returning id").Scan(&tokid)
		if err != nil {
			fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
			return
		}
	}

	// then i will insert it into yo mama
	query := `
		insert into tokens (service, token, refresh, userid, owner)
        values ($1, $2, $3, $4, $5)
	`
	_, err = db.Exec(query, "discord", tok.Token, tok.Refresh, user.ID, tokid)
	if err != nil {
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	// err = db.QueryRow("SELECT id, owner FROM tokens WHERE userid = $1", user.ID).Scan(&tokid, &owner)
	// if err != nil {
	// 	err = db.QueryRow("INSERT INTO tokens (service, token, refresh, userid) VALUES ($1, $2, $3, $4) RETURNING id",
	// 		"google", tok.Token, tok.Refresh, user.ID).Scan(&tokid)
	// 	if err != nil {
	// 		fmt.Fprintln(w, "Erreur lors de l'insertion du token:", err.Error())
	// 		return
	// 	}

	// 	err = db.QueryRow("INSERT INTO users (tokenid) VALUES ($1) RETURNING id", tokid).Scan(&owner)
	// 	if err != nil {
	// 		fmt.Fprintln(w, "Erreur lors de l'insertion de l'utilisateur:", err.Error())
	// 		return
	// 	}
	// 	db.Exec("UPDATE tokens SET owner = $1 WHERE id = $2", owner, tokid)
	// }

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

func sendEmail(gmailToken, recipientEmail, subject, body string) error {
	ctx := context.Background()
	token := &oauth2.Token{
		AccessToken: gmailToken,
	}
	tokenSource := oauth2.StaticTokenSource(token)
	srv, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return fmt.Errorf("Unable to create Gmail service: %v", err)
	}

	message := gmail.Message{}
	msgBody := "From: 'me'\r\n" +
		"To: " + recipientEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body
	message.Raw = base64.URLEncoding.EncodeToString([]byte(msgBody))

	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return fmt.Errorf("Unable to send email: %v", err)
	}

	return nil
}

func sendEmailNotification(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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
		Spices struct {
			Subject string `json:"subject"`
			Email   string `json:"email"`
			Body    string `json:"body"`
		} `json:"spices"`
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
	fmt.Println("Extracted userID:", userID)

	var gmailToken string
	err = db.QueryRow("SELECT token FROM tokens WHERE owner = $1 AND service = 'google'", userID).Scan(&gmailToken)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("No Google token found for user:", userID)
			fmt.Fprintf(w, "{ \"error\": \"No Google token found for user\" }\n")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Database error:", err.Error())
			fmt.Fprintf(w, "{ \"error\": \"Database error: %s\" }\n", err.Error())
		}
		return
	}

	if gmailToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("No Gmail token available for user:", userID)
		fmt.Fprintf(w, "{ \"error\": \"No Gmail token available\" }\n")
		return
	}

	spice := requestBody.Spices

	subject := spice.Subject
	recipientEmail := spice.Email
	body := spice.Body

	if recipientEmail == "" || subject == "" || body == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Missing email field (subject, email, or body) in spices")
		fmt.Fprintf(w, "{ \"error\": \"Missing one or more fields: subject, email, body\" }\n")
		return
	}

	err = sendEmail(gmailToken, recipientEmail, subject, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error sending email:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"Failed to send email\" }\n")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{ \"status\": \"Email sent successfully!\" }\n")
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
	Oauth  bool `json:"oauth"`
}

type Message struct {
	Bridge      int               `json:"bridge"`
	UserId      int               `json:"userid"`
	Ingredients map[string]string `json:"ingredients"`
}

func getRoutes(w http.ResponseWriter, req *http.Request) {
	var list = []InfoRoute{
		{
			Name: "sendEmailNotification",
			Type: "reaction",
			Desc: "send an email",
			Spices: []InfoSpice{
				{
					Name:  "email",
					Type:  "text",
					Title: "The email.",
				},
				{
					Name:  "subject",
					Type:  "text",
					Title: "The subject",
				},
				{
					Name:  "body",
					Type:  "text",
					Title: "The body you want to send.",
				},
			},
		},
	}
	var infos = Infos{
		Color:  "#4272db",
		Image:  "/assets/google.png",
		Routes: list,
		Oauth:  true,
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
	router.HandleFunc("/sendEmailNotification", miniproxy(sendEmailNotification, db)).Methods("POST")
	router.HandleFunc("/user", getUserInfo).Methods("GET")
	router.HandleFunc("/", getRoutes).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}
