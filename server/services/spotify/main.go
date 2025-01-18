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
        UserID  int    `json:"userid"`
        Spices  struct {
            Musique string `json:"musique"`
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
    trackName := requestBody.Spices.Musique
    fmt.Println("Extracted userID:", userID)
    fmt.Println("Requested trackName:", trackName)

    if trackName == "" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Println("Error: Track name is empty.")
        fmt.Fprintf(w, "{ \"error\": \"Track name is empty\" }\n")
        return
    }

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

    searchURL := "https://api.spotify.com/v1/search"
    query := fmt.Sprintf("%s", trackName)

    reqSearch, err := http.NewRequest("GET", searchURL+"?q="+url.QueryEscape(query)+"&type=track&limit=1", nil)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error creating search request:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }
    reqSearch.Header.Set("Authorization", "Bearer "+spotifyToken)

    client := &http.Client{}
    respSearch, err := client.Do(reqSearch)
    if err != nil {
        w.WriteHeader(http.StatusBadGateway)
        fmt.Println("Error searching track on Spotify:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }
    defer respSearch.Body.Close()

    fmt.Println("Search Response Status Code:", respSearch.StatusCode)
    bodyResp, err := io.ReadAll(respSearch.Body)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error reading search response body:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"Error reading search response body\" }\n")
        return
    }
    fmt.Println("Search Response Body:", string(bodyResp))

    if respSearch.StatusCode != http.StatusOK {
        w.WriteHeader(http.StatusNotFound)
        fmt.Println("Track not found. Status:", respSearch.StatusCode)
        fmt.Fprintf(w, "{ \"error\": \"Track not found\" }\n")
        return
    }

    var searchResult struct {
        Tracks struct {
            Items []struct {
                Uri string `json:"uri"`
            } `json:"items"`
        } `json:"tracks"`
    }
    if err := json.Unmarshal(bodyResp, &searchResult); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error unmarshalling search response:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling search response\" }\n")
        return
    }

    if len(searchResult.Tracks.Items) == 0 {
        w.WriteHeader(http.StatusNotFound)
        fmt.Println("No tracks found for the given name")
        fmt.Fprintf(w, "{ \"error\": \"No tracks found for the given name\" }\n")
        return
    }

    trackURI := searchResult.Tracks.Items[0].Uri
    fmt.Println("Track URI found:", trackURI)

    spotifyURL := "https://api.spotify.com/v1/me/player/play"
    playBody := fmt.Sprintf(`{"uris":["%s"]}`, trackURI)

    reqSpotify, err := http.NewRequest("PUT", spotifyURL, bytes.NewBufferString(playBody))
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error creating Spotify request:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }

    reqSpotify.Header.Set("Authorization", "Bearer "+spotifyToken)
    reqSpotify.Header.Set("Content-Type", "application/json")

    respSpotify, err := client.Do(reqSpotify)
    if err != nil {
        w.WriteHeader(http.StatusBadGateway)
        fmt.Println("Error playing music on Spotify:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }
    defer respSpotify.Body.Close()

    fmt.Println("Response Status Code:", respSpotify.StatusCode)
    bodyResp, err = io.ReadAll(respSpotify.Body)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error reading Spotify response body:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"Error reading Spotify response body\" }\n")
        return
    }
    fmt.Println("Spotify Response Body:", string(bodyResp))

    if respSpotify.StatusCode == http.StatusNoContent {
        fmt.Println("Music is now playing!")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "{ \"status\": \"Music is now playing!\" }\n")
    } else {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Failed to play music on Spotify. Status:", respSpotify.StatusCode)
        fmt.Fprintf(w, "{ \"error\": \"Failed to play music\" }\n")
    }
}

func pauseMusic(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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

    devicesURL := "https://api.spotify.com/v1/me/player/devices"
    reqDevices, err := http.NewRequest("GET", devicesURL, nil)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error creating devices request:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }
    reqDevices.Header.Set("Authorization", "Bearer "+spotifyToken)

    client := &http.Client{}
    respDevices, err := client.Do(reqDevices)
    if err != nil {
        w.WriteHeader(http.StatusBadGateway)
        fmt.Println("Error fetching devices:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }
    defer respDevices.Body.Close()

    bodyDevices, err := io.ReadAll(respDevices.Body)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error reading devices response body:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"Error reading devices response body\" }\n")
        return
    }
    fmt.Println("Devices Response Body:", string(bodyDevices))

    var devicesResult struct {
        Devices []struct {
            Active bool `json:"is_active"`
        } `json:"devices"`
    }

    if err := json.Unmarshal(bodyDevices, &devicesResult); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error unmarshalling devices response:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling devices response\" }\n")
        return
    }

    activeDeviceFound := false
    for _, device := range devicesResult.Devices {
        if device.Active {
            activeDeviceFound = true
            break
        }
    }

    if !activeDeviceFound {
        w.WriteHeader(http.StatusNotFound)
        fmt.Println("No active device found.")
        fmt.Fprintf(w, "{ \"error\": \"No active device found\" }\n")
        return
    }

    spotifyURL := "https://api.spotify.com/v1/me/player/pause"
    reqSpotify, err := http.NewRequest("PUT", spotifyURL, nil)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error creating Spotify request:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }

    reqSpotify.Header.Set("Authorization", "Bearer "+spotifyToken)

    respSpotify, err := client.Do(reqSpotify)
    if err != nil {
        w.WriteHeader(http.StatusBadGateway)
        fmt.Println("Error pausing music on Spotify:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }
    defer respSpotify.Body.Close()

    fmt.Println("Response Status Code:", respSpotify.StatusCode)
    bodyResp, err := io.ReadAll(respSpotify.Body)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error reading Spotify response body:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"Error reading Spotify response body\" }\n")
        return
    }
    fmt.Println("Spotify Response Body:", string(bodyResp))

    if respSpotify.StatusCode == http.StatusNoContent {
        fmt.Println("Music is paused!")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "{ \"status\": \"Music is paused!\" }\n")
    } else {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Failed to pause music on Spotify. Status:", respSpotify.StatusCode)
        fmt.Fprintf(w, "{ \"error\": \"Failed to pause music\" }\n")
    }
}

func resumeMusic(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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

    spotifyURL := "https://api.spotify.com/v1/me/player/play"
    reqSpotify, err := http.NewRequest("PUT", spotifyURL, nil)
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
        fmt.Println("Error resuming music on Spotify:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
        return
    }
    defer respSpotify.Body.Close()

    fmt.Println("Response Status Code:", respSpotify.StatusCode)
    bodyResp, err := io.ReadAll(respSpotify.Body)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Error reading response body:", err.Error())
        fmt.Fprintf(w, "{ \"error\": \"Error reading response body\" }\n")
        return
    }
    fmt.Println("Response Body:", string(bodyResp))

    if respSpotify.StatusCode == http.StatusNoContent {
        fmt.Println("Music is now playing!")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "{ \"status\": \"Music is now playing!\" }\n")
    } else {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Failed to resume music on Spotify. Status:", respSpotify.StatusCode)
        fmt.Fprintf(w, "{ \"error\": \"Failed to resume music\" }\n")
    }
}

func checkDeviceConnection(w http.ResponseWriter, req *http.Request, db *sql.DB) {
    bodyBytes, err := io.ReadAll(req.Body)
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

    var spotifyToken string
    err = db.QueryRow("SELECT token FROM tokens WHERE owner = $1 AND service = 'spotify'", userID).Scan(&spotifyToken)
    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "{ \"error\": \"No Spotify token found for user\" }\n")
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "{ \"error\": \"Database error: %s\" }\n", err.Error())
        }
        return
    }

    if spotifyToken == "" {
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprintf(w, "{ \"error\": \"No Spotify token available\" }\n")
        return
    }

    go func() {
        deviceURL := "https://api.spotify.com/v1/me/player/devices"
        reqDevices, err := http.NewRequest("GET", deviceURL, nil)
        if err != nil {
            fmt.Println("Error creating request:", err.Error())
            return
        }

        reqDevices.Header.Set("Authorization", "Bearer "+spotifyToken)

        client := &http.Client{}
        for {
            respDevices, err := client.Do(reqDevices)
            if err != nil {
                fmt.Println("Error fetching devices:", err.Error())
                return
            }
            defer respDevices.Body.Close()

            bodyResp, err := io.ReadAll(respDevices.Body)
            if err != nil {
                fmt.Println("Error reading response body:", err.Error())
                return
            }

            if respDevices.StatusCode != http.StatusOK {
                fmt.Println("Failed to get devices:", respDevices.StatusCode)
                return
            }

            var deviceResponse struct {
                Devices []struct {
                    ID       string `json:"id"`
                    IsActive bool   `json:"is_active"`
                    Name     string `json:"name"`
                } `json:"devices"`
            }

            if err := json.Unmarshal(bodyResp, &deviceResponse); err != nil {
                fmt.Println("Error unmarshalling response:", err.Error())
                return
            }

            var activeDevice *struct {
                ID       string `json:"id"`
                IsActive bool   `json:"is_active"`
                Name     string `json:"name"`
            }

            for _, device := range deviceResponse.Devices {
                if device.IsActive {
                    activeDevice = &device
                    break
                }
            }

            if activeDevice != nil {
                url := fmt.Sprintf("http://backend:%d/api/orchestrator", bridgeID)
                requestBody := map[string]interface{}{
                    "bridge":     bridgeID,
                    "userid":     userID,
                    "ingredients": map[string]interface{}{},
                }
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

                fmt.Println("Device is connected and active:", activeDevice.Name)
                return
            }

            time.Sleep(1 * time.Second)
        }
    }()

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "{ \"status\": \"Checking devices...\" }\n")
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
	fmt.Println("Spotify microservice container is running !")
	router := mux.NewRouter()
	godotenv.Load(".env")

	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	router.HandleFunc("/playMusic", miniproxy(playMusic, db)).Methods("POST")
	router.HandleFunc("/pauseMusic", miniproxy(pauseMusic, db)).Methods("POST")
	router.HandleFunc("/resumeMusic", miniproxy(resumeMusic, db)).Methods("POST")
	router.HandleFunc("/checkDeviceConnection", miniproxy(checkDeviceConnection, db)).Methods("POST")
	router.HandleFunc("/user", getUserInfo).Methods("GET")
	router.HandleFunc("/", getRoutes).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}
