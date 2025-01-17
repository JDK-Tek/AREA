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

type TeamsContent struct {
	TeamID   string `json:"team_id"`
	ChannelID string `json:"channel_id"`
	Message  string `json:"message"`
}

type Objects struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

type Content struct {
	Dishes Objects `json:"spices"`
}

type EmailContent struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func getOAUTHLink(w http.ResponseWriter, req *http.Request) {
	str := "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?"
	x := url.QueryEscape(os.Getenv("REDIRECT"))
	str += "client_id=" + os.Getenv("OUTLOOK_CLIENT_ID")
	str += "&response_type=code"
	str += "&redirect_uri=" + x
	str += "&scope=User.Read"
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

	clientid := os.Getenv("OUTLOOK_CLIENT_ID")
	clientsecret := os.Getenv("OUTLOOK_CLIENT_SECRET")
	data := url.Values{}
	err := json.NewDecoder(req.Body).Decode(&res)
	if err != nil {
		fmt.Fprintln(w, "decodeeee", err.Error())
		return
	}
	data.Set("client_id", clientid)
	data.Set("client_secret", clientsecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", res.Code)
	data.Set("redirect_uri", os.Getenv("REDIRECT"))
	rep, err := http.PostForm(API_OAUTH_OUTLOOK, data)
	// fmt.Fprintln(w, "tmp = ", string(body))
	// return
	if err != nil {
		fmt.Fprintln(w, "postform", err.Error())
		return
	}
	defer rep.Body.Close()
	body, err := io.ReadAll(rep.Body)
	if err != nil {
		fmt.Fprintln(w, "read body", err.Error())
		return
	}
	fmt.Println("token value = ", rep.Body)
	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Fprintln(w, "unmarshal json", err.Error())
		return
	}
	fmt.Println(responseData)
	tok.Token = responseData["access_token"].(string)
	tok.Refresh = "foo bar"

	req, err = http.NewRequest("GET", API_USER_OUTLOOK, nil)
	if err != nil {
		fmt.Fprintln(w, "request error", err.Error())
		return
	}
	req.Header.Set("Authorization", "Bearer " +tok.Token)
	client := &http.Client{}
	rep, err = client.Do(req)
	if err != nil {
		fmt.Fprintln(w, "client do", err.Error())
		return
	}
	defer rep.Body.Close()
	err = json.NewDecoder(rep.Body).Decode(&user)
	if err != nil {
		fmt.Fprintln(w, "decodeiiiiii", err.Error())
		return
	}

	if tok.Token == "" || tok.Refresh == "" {
		fmt.Fprintln(w, "error: token is empty")
		return
	}
	err = db.QueryRow("select id, owner from tokens where userid = $1", user.ID).Scan(&tokid, &owner)
	if err != nil {
		err = db.QueryRow("insert into tokens (service, token, refresh, userid) values ($1, $2, $3, $4) returning id",
			"outlook",
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
	fmt.Println("Sucess login with outlook, token = ", tokenStr)
	fmt.Fprintf(w, "{\"token\": \"%s\"}\n", tokenStr)
}

func sendTeamsMessage(w http.ResponseWriter, req *http.Request) {
	var teamsContent TeamsContent
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&teamsContent)
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

	messageData := map[string]interface{}{
		"body": map[string]interface{}{
			"content": teamsContent.Message,
		},
	}
	messageBytes, err := json.Marshal(messageData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	teamsMessageURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages", teamsContent.TeamID, teamsContent.ChannelID)

	reqTeams, err := http.NewRequest("POST", teamsMessageURL, bytes.NewBuffer(messageBytes))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqTeams.Header.Set("Authorization", "Bearer "+token)
	reqTeams.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqTeams)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{ \"status\": \"Message sent successfully\" }\n")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"Failed to send message\" }\n")
	}
}

func sendEmail(w http.ResponseWriter, req *http.Request, db *sql.DB) {
    fmt.Println("Headers received:", req.Header)

    bodyBytes, err := io.ReadAll(req.Body)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }

    fmt.Println("Request Body:", string(bodyBytes))

    req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

    var requestBody struct {
        UserID int `json:"userid"`
        Spices struct {
            Message string `json:"message"`
            To      string `json:"to"`
            Subject string `json:"subject"`
        } `json:"spices"`
    }

    decoder := json.NewDecoder(req.Body)
    err = decoder.Decode(&requestBody)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }

    userID := requestBody.UserID

    var outlookToken string
    err = db.QueryRow("SELECT token FROM tokens WHERE owner = $1 AND service = 'outlook'", userID).Scan(&outlookToken)
    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "{ \"error\": \"No Outlook token found for user\" }\n")
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "{ \"error\": \"Database error: %s\" }\n", err.Error())
        }
        return
    }

    if outlookToken == "" {
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "{ \"error\": \"No Outlook token available\" }\n")
        return
    }

    emailContent := requestBody.Spices
    if emailContent.To == "" || emailContent.Subject == "" || emailContent.Message == "" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "{ \"error\": \"Missing email details\" }\n")
        return
    }

    emailData := map[string]interface{}{
        "message": map[string]interface{}{
            "subject": emailContent.Subject,
            "body": map[string]interface{}{
                "contentType": "Text",
                "content":     emailContent.Message,
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

    reqEmail.Header.Set("Authorization", "Bearer "+outlookToken)
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

func checkEmail(w http.ResponseWriter, req *http.Request) {
	token := os.Getenv("OUTLOOK_TOKEN")
	if token == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"token is missing\" }\n")
		return
	}

	emailAPI := "https://graph.microsoft.com/v1.0/me/messages?$top=5&$orderby=receivedDateTime desc"

	reqEmail, err := http.NewRequest("GET", emailAPI, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqEmail.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(reqEmail)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"Failed to fetch emails\" }\n")
		return
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	emails, ok := response["value"].([]interface{})
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"Failed to parse emails\" }\n")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	emailList := make([]map[string]interface{}, len(emails))
	for i, email := range emails {
		emailData, _ := email.(map[string]interface{})
		emailList[i] = map[string]interface{}{
			"subject": emailData["subject"],
			"from":    emailData["from"],
			"received": emailData["receivedDateTime"],
		}
	}
	json.NewEncoder(w).Encode(emailList)
}

func checkTeamsMessages(w http.ResponseWriter, req *http.Request) {
	var teamsContent TeamsContent
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&teamsContent)
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

	teamsMessageURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages?$top=5&$orderby=createdDateTime desc", teamsContent.TeamID, teamsContent.ChannelID)

	reqTeams, err := http.NewRequest("GET", teamsMessageURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqTeams.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(reqTeams)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"Failed to fetch Teams messages\" }\n")
		return
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	messages, ok := response["value"].([]interface{})
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"Failed to parse messages\" }\n")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	messageList := make([]map[string]interface{}, len(messages))
	for i, message := range messages {
		messageData, _ := message.(map[string]interface{})
		messageList[i] = map[string]interface{}{
			"message":   messageData["body"],
			"from":      messageData["from"],
			"createdAt": messageData["createdDateTime"],
		}
	}
	json.NewEncoder(w).Encode(messageList)
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
			Name: "sendEmail",
			Type: "reaction",
			Desc: "Sends an email.",
			Spices: []InfoSpice{
				{
					Name: "to",
					Type: "text",
					Title: "email to be sending",
				},
				{
					Name: "subject",
					Type: "text",
					Title: "The subject",
				},
				{
					Name: "body",
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
	fmt.Println("outlook microservice container is running !")
	router := mux.NewRouter()
	godotenv.Load(".env")
	//router.HandleFunc("/send", doSomeSend).Methods("POST")
	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	router.HandleFunc("/sendEmail", miniproxy(sendEmail, db)).Methods("POST")
	router.HandleFunc("/sendTeamsMessage", sendTeamsMessage).Methods("POST")
	router.HandleFunc("/checkEmail", checkEmail).Methods("GET")
	router.HandleFunc("/checkTeamsMessages", checkTeamsMessages).Methods("POST")
	router.HandleFunc("/", getRoutes).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}