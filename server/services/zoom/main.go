package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	ZOOM_AUTH_URL   = "https://zoom.us/oauth/authorize"
	ZOOM_TOKEN_URL  = "https://zoom.us/oauth/token"
	ZOOM_API_ME_URL = "https://api.zoom.us/v2/users/me"
)

const PERMISSIONS = 8
const EXPIRATION = 60 * 30

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
	Oauth  bool        `json:"oauth"`
}

func getUserInfo(w http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Token d'accès manquant", http.StatusUnauthorized)
		return
	}

	reqUser, err := http.NewRequest("GET", ZOOM_API_ME_URL, nil)
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

func getOAUTHLink(w http.ResponseWriter, req *http.Request) {
    fmt.Println("test test")
    str := "https://zoom.us/oauth/authorize?"

    redirectURI := url.QueryEscape(os.Getenv("REDIRECT"))
    fmt.Println("Redirect URI = ", redirectURI)

    scopes := "meeting:read:list_meetings meeting:read:list_meetings:admin"

    str += "response_type=code"
    str += "&state=some-state-value"
    str += "&client_id=" + os.Getenv("ZOOM_CLIENT_ID")
    str += "&redirect_uri=" + redirectURI
    str += "&scope=" + url.QueryEscape(scopes)

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, str)
}


func miniproxy(f func(http.ResponseWriter, *http.Request, *sql.DB), c *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(a http.ResponseWriter, b *http.Request) {
		f(a, b, c)
	}
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

    clientID := os.Getenv("ZOOM_CLIENT_ID")
    clientSecret := os.Getenv("ZOOM_CLIENT_SECRET")
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
   
    rep, err := http.PostForm("https://zoom.us/oauth/token", data)
    if err != nil {
        fmt.Fprintln(w, "Erreur lors de l'échange du code:", err.Error())
        return
    }
    defer rep.Body.Close()
   
    body, err := io.ReadAll(rep.Body)
	fmt.Println("body = ", string(body))
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

	fmt.Println("token = ", tok.Token)
	fmt.Println("refresh = ", tok.Refresh)

    req, err = http.NewRequest("GET", "https://api.zoom.us/v2/users/me", nil)
    if err != nil {
        fmt.Fprintln(w, "Erreur lors de la création de la requête utilisateur:", err.Error())
        return
    }
    req.Header.Set("Authorization", "Bearer "+tok.Token)
   
    client := &http.Client{}
    rep, err = client.Do(req)
    if err != nil {
        fmt.Fprintln(w, "Erreur lors de l'appel à l'API Zoom:", err.Error())
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
            "zoom",
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
   
    fmt.Println("Succès de l'authentification avec Zoom, token =", tokenStr)
    fmt.Fprintf(w, "{\"token\": \"%s\"}\n", tokenStr)
}

type Message struct {
    Bridge int `json:"bridge"`
    UserId int `json:"userid"`
    Ingredients map[string]string `json:"ingredients"`
}

func checkDeviceConnection(w http.ResponseWriter, req *http.Request, db *sql.DB) {
    bodyBytes, err := io.ReadAll(req.Body)
    backendPort := os.Getenv("BACKEND_PORT")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }

    req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

    var requestBody struct {
        UserID int `json:"userid"`
        Bridge int `json:"bridge"`
        Spices struct{} `json:"spices"`
    }

    decoder := json.NewDecoder(req.Body)
    err = decoder.Decode(&requestBody)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }

    userID := requestBody.UserID
    bridgeID := requestBody.Bridge

    fmt.Println("UserID:", userID)
    fmt.Println("BridgeID:", bridgeID)

    var zoomToken string
    err = db.QueryRow("SELECT token FROM tokens WHERE owner = $1 AND service = 'zoom'", userID).Scan(&zoomToken)
    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "{ \"error\": \"No Zoom token found for user\" }\n")
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "{ \"error\": \"Database error: %s\" }\n", err.Error())
        }
        return
    }

    if zoomToken == "" {
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "{ \"error\": \"No Zoom token available\" }\n")
        return
    }

    fmt.Println("Zoom token retrieved successfully")

    go func() {
        client := &http.Client{}

        userInfoURL := "https://api.zoom.us/v2/users/me"
        reqUserInfo, err := http.NewRequest("GET", userInfoURL, nil)
        if err != nil {
            fmt.Println("Error creating request:", err.Error())
            return
        }
        reqUserInfo.Header.Set("Authorization", "Bearer "+zoomToken)

        fmt.Println("Fetching Zoom user info...")
        respUserInfo, err := client.Do(reqUserInfo)
        if err != nil {
            fmt.Println("Error fetching user info:", err.Error())
            return
        }
        defer respUserInfo.Body.Close()

        bodyResp, err := io.ReadAll(respUserInfo.Body)
        if err != nil {
            fmt.Println("Error reading response body:", err.Error())
            return
        }

        if respUserInfo.StatusCode != http.StatusOK {
            fmt.Println("Failed to get user info:", respUserInfo.StatusCode)
            fmt.Println("Response body:", string(bodyResp))
            return
        }

        var userInfoResponse struct {
            ID    string `json:"id"`
            Email string `json:"email"`
        }

        if err := json.Unmarshal(bodyResp, &userInfoResponse); err != nil {
            fmt.Println("Error unmarshalling response:", err.Error())
            return
        }

        zoomUserID := userInfoResponse.ID
        fmt.Println("Zoom UserID:", zoomUserID)

        userMeetingsURL := fmt.Sprintf("https://api.zoom.us/v2/users/%s/meetings", zoomUserID)
        fmt.Println("Fetching Zoom meetings for user...")

        for {
            reqMeetings, err := http.NewRequest("GET", userMeetingsURL, nil)
            if err != nil {
                fmt.Println("Error creating request for meetings:", err.Error())
                return
            }
            reqMeetings.Header.Set("Authorization", "Bearer "+zoomToken)

            respMeetings, err := client.Do(reqMeetings)
            if err != nil {
                fmt.Println("Error fetching meetings:", err.Error())
                return
            }
            defer respMeetings.Body.Close()

            bodyResp, err := io.ReadAll(respMeetings.Body)
            if err != nil {
                fmt.Println("Error reading response body:", err.Error())
                return
            }

            if respMeetings.StatusCode != http.StatusOK {
                fmt.Println("Failed to get meetings:", respMeetings.StatusCode)
                fmt.Println("Response body:", string(bodyResp)) // Added to show error details
                return
            }

            var meetingsResponse struct {
                Meetings []struct {
                    ID        string `json:"id"`
                    Topic     string `json:"topic"`
                    StartTime string `json:"start_time"`
                    Status    string `json:"status"`
                } `json:"meetings"`
            }

            if err := json.Unmarshal(bodyResp, &meetingsResponse); err != nil {
                fmt.Println("Error unmarshalling meetings response:", err.Error())
                return
            }

            var activeMeeting *struct {
                ID        string `json:"id"`
                Topic     string `json:"topic"`
                StartTime string `json:"start_time"`
                Status    string `json:"status"`
            }

            for _, meeting := range meetingsResponse.Meetings {
                if meeting.Status == "started" {
                    activeMeeting = &meeting
                    break
                }
            }

            if activeMeeting != nil {
                fmt.Println("User is in an active meeting:", activeMeeting.Topic)

                url := fmt.Sprintf("http://backend:%s/api/orchestrator", backendPort)
                var requestBody Message

                requestBody.Bridge = bridgeID
                requestBody.UserId = userID
                requestBody.Ingredients = make(map[string]string)
                requestBody.Ingredients["zoom_status"] = "in_meeting"

                jsonData, err := json.Marshal(requestBody)
                if err != nil {
                    fmt.Println("Error marshaling JSON:", err.Error())
                    return
                }

                reqPut, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
                if err != nil {
                    fmt.Println("Error creating PUT request:", err.Error())
                    return
                }

                reqPut.Header.Set("Content-Type", "application/json")

                respPut, err := client.Do(reqPut)
                if err != nil {
                    fmt.Println("Error sending PUT request:", err.Error())
                    return
                }
                defer respPut.Body.Close()

                if respPut.StatusCode != http.StatusOK {
                    fmt.Println("Failed to send PUT request:", respPut.StatusCode)
                    return
                }

                fmt.Println("User meeting information updated successfully")
                return
            }

            fmt.Println("User is not in a meeting, retrying...")
            time.Sleep(1 * time.Second)
        }
    }()

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "{ \"status\": \"ok !\" }\n")
}

func getRoutes(w http.ResponseWriter, req *http.Request) {
	var list = []InfoRoute{
		{
			Name: "checkDeviceConnection",
			Type: "action",
			Desc: "Check if the user is in meating !",
			Spices: []InfoSpice{},
		},
	}
	var infos = Infos{
		Color:  "#1DB954",
		Image:  "/assets/spotify.webp",
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

func connectToDatabase() (*sql.DB, error) {
	fmt.Println("p")
	dbPassword := os.Getenv("DB_PASSWORD")
	fmt.Println("a")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD not found")
	}
	fmt.Println("b")
	dbUser := os.Getenv("DB_USER")
	fmt.Println("c")
	if dbUser == "" {
		log.Fatal("DB_USER not found")
	}
	fmt.Println("d")
	fmt.Println("z")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal("DB_HOST not found")
	}
	fmt.Println("e")
	dbName := os.Getenv("DB_NAME")
	fmt.Println("zizi")
	if dbName == "" {
		log.Fatal("DB_NAME not found")
	}
	fmt.Println("co")
	dbPort := os.Getenv("DB_PORT")
	fmt.Println("u")
	if dbPort == "" {
		log.Fatal("DB_PORT not found")
	}
	fmt.Println("wawa")
	connectStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)
	fmt.Println("coucou")
	return sql.Open("postgres", connectStr)
}

func main() {
	godotenv.Load(".env")
	godotenv.Load("/usr/mount.d/.env")
	db, err := connectToDatabase()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(84)
	}
    fmt.Println("Zoom microservice container is running !")
    router := mux.NewRouter()
	
    router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
    router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	router.HandleFunc("/checkDeviceConnection", miniproxy(checkDeviceConnection, db)).Methods("POST")
    router.HandleFunc("/user", getUserInfo).Methods("POST")
    router.HandleFunc("/", getRoutes).Methods("GET")
    log.Fatal(http.ListenAndServe(":80", router))
}
