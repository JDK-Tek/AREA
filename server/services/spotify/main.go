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
	"strings"
	"time"

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
		"user-follow-modify " +
		"user-library-modify"

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
		UserID int `json:"userid"`
		Spices struct {
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

func likeCurrentMusic(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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

	spotifyURL := "https://api.spotify.com/v1/me/player/currently-playing"
	reqSpotify, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error creating Spotify request:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqSpotify.Header.Set("Authorization", "Bearer "+spotifyToken)
	client := &http.Client{}
	respSpotify, err := client.Do(reqSpotify)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Println("Error fetching currently playing music:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer respSpotify.Body.Close()

	bodyResp, err := io.ReadAll(respSpotify.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error reading Spotify response body:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"Error reading Spotify response body\" }\n")
		return
	}
	fmt.Println("Spotify Response Body:", string(bodyResp))

	if respSpotify.StatusCode == http.StatusOK {
		var playbackResponse struct {
			Item struct {
				Name    string `json:"name"`
				ID      string `json:"id"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
			} `json:"item"`
		}

		if err := json.Unmarshal(bodyResp, &playbackResponse); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error unmarshalling playback response:", err.Error())
			fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling playback response\" }\n")
			return
		}

		if playbackResponse.Item.Name != "" {
			trackName := playbackResponse.Item.Name
			trackID := playbackResponse.Item.ID // L'ID du morceau
			artistName := playbackResponse.Item.Artists[0].Name

			fmt.Printf("Currently playing: %s by %s\n", trackName, artistName)

			// Add the track to the user's library (like it)
			likeTrackURL := "https://api.spotify.com/v1/me/tracks?ids=" + trackID
			reqLike, err := http.NewRequest("PUT", likeTrackURL, nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error creating like request:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
				return
			}

			reqLike.Header.Set("Authorization", "Bearer "+spotifyToken)

			respLike, err := client.Do(reqLike)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				fmt.Println("Error liking track:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
				return
			}
			defer respLike.Body.Close()

			if respLike.StatusCode == http.StatusOK {
				fmt.Println("Track has been liked!")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "{ \"status\": \"Track has been liked!\" }\n")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Failed to like track. Status:", respLike.StatusCode)
				fmt.Fprintf(w, "{ \"error\": \"Failed to like track\" }\n")
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("No track is currently playing.")
			fmt.Fprintf(w, "{ \"error\": \"No track is currently playing\" }\n")
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed to fetch currently playing track. Status:", respSpotify.StatusCode)
		fmt.Fprintf(w, "{ \"error\": \"Failed to fetch currently playing track\" }\n")
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
	backendPort := os.Getenv("BACKEND_PORT")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	var requestBody struct {
		UserID int      `json:"userid"`
		Bridge int      `json:"bridge"`
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

	fmt.Println("userID = ", userID)
	fmt.Println("bridge = ", bridgeID)
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
				url := fmt.Sprintf("http://backend:%s/api/orchestrator", backendPort)
				var requestBody Message

				requestBody.Bridge = bridgeID
				requestBody.UserId = userID
				requestBody.Ingredients = make(map[string]string)
				requestBody.Ingredients["device"] = string(activeDevice.Name)
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
	fmt.Fprintf(w, "{ \"status\": \"ok !\" }\n")
}

func checkSongRunning(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	bodyBytes, err := io.ReadAll(req.Body)
	backendPort := os.Getenv("BACKEND_PORT")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	var requestBody struct {
		UserID int      `json:"userid"`
		Bridge int      `json:"bridge"`
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

	fmt.Println("userID = ", userID)
	fmt.Println("bridge = ", bridgeID)

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
				playbackURL := "https://api.spotify.com/v1/me/player/currently-playing"
				reqPlayback, err := http.NewRequest("GET", playbackURL, nil)
				if err != nil {
					fmt.Println("Error creating playback request:", err.Error())
					return
				}

				reqPlayback.Header.Set("Authorization", "Bearer "+spotifyToken)

				respPlayback, err := client.Do(reqPlayback)
				if err != nil {
					fmt.Println("Error fetching playback info:", err.Error())
					return
				}
				defer respPlayback.Body.Close()

				bodyPlayback, err := io.ReadAll(respPlayback.Body)
				if err != nil {
					fmt.Println("Error reading playback response body:", err.Error())
					return
				}

				if respPlayback.StatusCode != http.StatusOK {
					fmt.Println("No track currently playing or error retrieving track.")
				} else {
					var playbackResponse struct {
						Item struct {
							Name    string `json:"name"`
							Artists []struct {
								Name string `json:"name"`
								ID   string `json:"id"`
							} `json:"artists"`
						} `json:"item"`
					}

					if err := json.Unmarshal(bodyPlayback, &playbackResponse); err != nil {
						fmt.Println("Error unmarshalling playback response:", err.Error())
						return
					}

					if playbackResponse.Item.Name != "" {
						trackName := playbackResponse.Item.Name
						artistName := playbackResponse.Item.Artists[0].Name
						artistID := playbackResponse.Item.Artists[0].ID

						artistURL := fmt.Sprintf("https://api.spotify.com/v1/artists/%s", artistID)
						reqArtist, err := http.NewRequest("GET", artistURL, nil)
						if err != nil {
							fmt.Println("Error creating artist request:", err.Error())
							return
						}

						reqArtist.Header.Set("Authorization", "Bearer "+spotifyToken)

						respArtist, err := client.Do(reqArtist)
						if err != nil {
							fmt.Println("Error fetching artist info:", err.Error())
							return
						}
						defer respArtist.Body.Close()

						bodyArtist, err := io.ReadAll(respArtist.Body)
						if err != nil {
							fmt.Println("Error reading artist response body:", err.Error())
							return
						}

						if respArtist.StatusCode != http.StatusOK {
							fmt.Println("Failed to get artist genres:", respArtist.StatusCode)
						} else {
							var artistResponse struct {
								Genres []string `json:"genres"`
							}

							if err := json.Unmarshal(bodyArtist, &artistResponse); err != nil {
								fmt.Println("Error unmarshalling artist response:", err.Error())
								return
							}

							genre := ""
							if len(artistResponse.Genres) > 0 {
								genre = artistResponse.Genres[0]
							}

							fmt.Printf("Currently playing: %s by %s, Genre: %s\n", trackName, artistName, genre)

							url := fmt.Sprintf("http://backend:%s/api/orchestrator", backendPort)
							var requestBody Message

							requestBody.Bridge = bridgeID
							requestBody.UserId = userID
							requestBody.Ingredients = make(map[string]string)
							requestBody.Ingredients["device"] = string(activeDevice.Name)
							requestBody.Ingredients["track"] = trackName
							requestBody.Ingredients["artist"] = artistName
							requestBody.Ingredients["genre"] = genre

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

							fmt.Println("Successfully sent the active track info to backend!")
							return
						}
					}
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{ \"status\": \"ok !\" }\n")
}

func checkPodcastRunning(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	bodyBytes, err := io.ReadAll(req.Body)
	backendPort := os.Getenv("BACKEND_PORT")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	var requestBody struct {
		UserID int      `json:"userid"`
		Bridge int      `json:"bridge"`
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
				playbackURL := "https://api.spotify.com/v1/me/player/currently-playing"
				reqPlayback, err := http.NewRequest("GET", playbackURL, nil)
				if err != nil {
					fmt.Println("Error creating playback request:", err.Error())
					return
				}

				reqPlayback.Header.Set("Authorization", "Bearer "+spotifyToken)

				respPlayback, err := client.Do(reqPlayback)
				if err != nil {
					fmt.Println("Error fetching playback info:", err.Error())
					return
				}
				defer respPlayback.Body.Close()

				bodyPlayback, err := io.ReadAll(respPlayback.Body)
				if err != nil {
					fmt.Println("Error reading playback response body:", err.Error())
					return
				}

				if respPlayback.StatusCode != http.StatusOK {
					fmt.Println("No track or podcast currently playing or error retrieving track.")
				} else {
					var playbackResponse struct {
						Item struct {
							Name string `json:"name"`
							Show struct {
								Name string `json:"name"`
								ID   string `json:"id"`
							} `json:"show"`
							Type string `json:"type"`
						} `json:"item"`
					}

					if err := json.Unmarshal(bodyPlayback, &playbackResponse); err != nil {
						fmt.Println("Error unmarshalling playback response:", err.Error())
						return
					}

					if playbackResponse.Item.Type == "episode" {
						podcastName := playbackResponse.Item.Show.Name
						episodeName := playbackResponse.Item.Name
						podcastID := playbackResponse.Item.Show.ID

						fmt.Printf("Currently playing podcast: %s - Episode: %s\n", podcastName, episodeName)

						episodesURL := fmt.Sprintf("https://api.spotify.com/v1/shows/%s/episodes", podcastID)
						reqEpisodes, err := http.NewRequest("GET", episodesURL, nil)
						if err != nil {
							fmt.Println("Error creating episodes request:", err.Error())
							return
						}

						reqEpisodes.Header.Set("Authorization", "Bearer "+spotifyToken)

						respEpisodes, err := client.Do(reqEpisodes)
						if err != nil {
							fmt.Println("Error fetching episodes:", err.Error())
							return
						}
						defer respEpisodes.Body.Close()

						bodyEpisodes, err := io.ReadAll(respEpisodes.Body)
						if err != nil {
							fmt.Println("Error reading episodes response body:", err.Error())
							return
						}

						if respEpisodes.StatusCode != http.StatusOK {
							fmt.Println("Failed to get episodes:", respEpisodes.StatusCode)
						} else {
							var episodesResponse struct {
								Items []struct {
									Name string `json:"name"`
									ID   string `json:"id"`
								} `json:"items"`
							}

							if err := json.Unmarshal(bodyEpisodes, &episodesResponse); err != nil {
								fmt.Println("Error unmarshalling episodes response:", err.Error())
								return
							}

							fmt.Println("Episodes of the podcast:")
							for _, episode := range episodesResponse.Items {
								fmt.Println("- " + episode.Name)
							}

							followURL := fmt.Sprintf("https://api.spotify.com/v1/me/following?type=show&ids=%s", podcastID)
							reqFollow, err := http.NewRequest("PUT", followURL, nil)
							if err != nil {
								fmt.Println("Error creating follow request:", err.Error())
								return
							}

							reqFollow.Header.Set("Authorization", "Bearer "+spotifyToken)

							respFollow, err := client.Do(reqFollow)
							if err != nil {
								fmt.Println("Error following podcast:", err.Error())
								return
							}
							defer respFollow.Body.Close()

							if respFollow.StatusCode != http.StatusOK {
								fmt.Println("Failed to follow podcast:", respFollow.StatusCode)
							} else {
								url := fmt.Sprintf("http://backend:%s/api/orchestrator", backendPort)
								var requestBody Message

								requestBody.Bridge = bridgeID
								requestBody.UserId = userID
								requestBody.Ingredients = make(map[string]string)
								requestBody.Ingredients["device"] = string(activeDevice.Name)
								requestBody.Ingredients["podcast"] = podcastName
								requestBody.Ingredients["episode"] = episodeName

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

								fmt.Println("Successfully followed the podcast and sent the information to backend!")
							}
						}
					}
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{ \"status\": \"ok !\" }\n")
}

func unlikeCurrentMusic(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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

	spotifyURL := "https://api.spotify.com/v1/me/player/currently-playing"
	reqSpotify, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error creating Spotify request:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqSpotify.Header.Set("Authorization", "Bearer "+spotifyToken)
	client := &http.Client{}
	respSpotify, err := client.Do(reqSpotify)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Println("Error fetching currently playing music:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer respSpotify.Body.Close()

	bodyResp, err := io.ReadAll(respSpotify.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error reading Spotify response body:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"Error reading Spotify response body\" }\n")
		return
	}
	fmt.Println("Spotify Response Body:", string(bodyResp))

	if respSpotify.StatusCode == http.StatusOK {
		var playbackResponse struct {
			Item struct {
				Name    string `json:"name"`
				ID      string `json:"id"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
			} `json:"item"`
		}

		if err := json.Unmarshal(bodyResp, &playbackResponse); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error unmarshalling playback response:", err.Error())
			fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling playback response\" }\n")
			return
		}

		if playbackResponse.Item.Name != "" {
			trackID := playbackResponse.Item.ID
			artistName := playbackResponse.Item.Artists[0].Name

			fmt.Printf("Currently playing: %s by %s\n", playbackResponse.Item.Name, artistName)

			unlikeTrackURL := "https://api.spotify.com/v1/me/tracks?ids=" + trackID
			reqUnlike, err := http.NewRequest("DELETE", unlikeTrackURL, nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error creating unlike request:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
				return
			}

			reqUnlike.Header.Set("Authorization", "Bearer "+spotifyToken)

			respUnlike, err := client.Do(reqUnlike)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				fmt.Println("Error unliking track:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
				return
			}
			defer respUnlike.Body.Close()

			if respUnlike.StatusCode == http.StatusOK {
				fmt.Println("Track has been unliked!")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "{ \"status\": \"Track has been unliked!\" }\n")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Failed to unlike track. Status:", respUnlike.StatusCode)
				fmt.Fprintf(w, "{ \"error\": \"Failed to unlike track\" }\n")
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("No track is currently playing.")
			fmt.Fprintf(w, "{ \"error\": \"No track is currently playing\" }\n")
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed to fetch currently playing track. Status:", respSpotify.StatusCode)
		fmt.Fprintf(w, "{ \"error\": \"Failed to fetch currently playing track\" }\n")
	}
}

func nextMusic(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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

	spotifyURL := "https://api.spotify.com/v1/me/player/next"
	reqSpotify, err := http.NewRequest("POST", spotifyURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error creating Spotify request:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqSpotify.Header.Set("Authorization", "Bearer "+spotifyToken)

	client := &http.Client{}
	respSpotify, err := client.Do(reqSpotify)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Println("Error moving to next track:", err.Error())
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

	if respSpotify.StatusCode == http.StatusOK {
		fmt.Println("Successfully skipped to the next track!")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{ \"status\": \"Successfully skipped to the next track!\" }\n")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed to skip to next track. Status:", respSpotify.StatusCode)
		fmt.Fprintf(w, "{ \"error\": \"Failed to skip to next track\" }\n")
	}
}

func previousMusic(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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

	spotifyURL := "https://api.spotify.com/v1/me/player/previous"
	reqSpotify, err := http.NewRequest("POST", spotifyURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error creating Spotify request:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqSpotify.Header.Set("Authorization", "Bearer "+spotifyToken)

	client := &http.Client{}
	respSpotify, err := client.Do(reqSpotify)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Println("Error moving to previous track:", err.Error())
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

	if respSpotify.StatusCode == http.StatusOK {
		fmt.Println("Successfully skipped to the previous track!")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{ \"status\": \"Successfully skipped to the previous track!\" }\n")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed to skip to previous track. Status:", respSpotify.StatusCode)
		fmt.Fprintf(w, "{ \"error\": \"Failed to skip to previous track\" }\n")
	}
}

func removeFromPlaylistIfPresent(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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

	spotifyURL := "https://api.spotify.com/v1/me/player/currently-playing"
	reqSpotify, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error creating Spotify request:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqSpotify.Header.Set("Authorization", "Bearer "+spotifyToken)
	client := &http.Client{}
	respSpotify, err := client.Do(reqSpotify)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Println("Error fetching currently playing music:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer respSpotify.Body.Close()

	bodyResp, err := io.ReadAll(respSpotify.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error reading Spotify response body:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"Error reading Spotify response body\" }\n")
		return
	}
	fmt.Println("Spotify Response Body:", string(bodyResp))

	if respSpotify.StatusCode == http.StatusOK {
		var playbackResponse struct {
			Item struct {
				Name    string `json:"name"`
				ID      string `json:"id"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
			} `json:"item"`
		}

		if err := json.Unmarshal(bodyResp, &playbackResponse); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error unmarshalling playback response:", err.Error())
			fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling playback response\" }\n")
			return
		}

		if playbackResponse.Item.Name != "" {
			trackID := playbackResponse.Item.ID
			trackName := playbackResponse.Item.Name
			artistName := playbackResponse.Item.Artists[0].Name

			fmt.Printf("Currently playing: %s by %s\n", trackName, artistName)

			playlistsURL := "https://api.spotify.com/v1/me/playlists"
			reqPlaylists, err := http.NewRequest("GET", playlistsURL, nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error creating playlists request:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
				return
			}

			reqPlaylists.Header.Set("Authorization", "Bearer "+spotifyToken)
			respPlaylists, err := client.Do(reqPlaylists)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				fmt.Println("Error fetching playlists:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
				return
			}
			defer respPlaylists.Body.Close()

			bodyPlaylists, err := io.ReadAll(respPlaylists.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error reading playlists response body:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"Error reading playlists response body\" }\n")
				return
			}

			var playlistsResponse struct {
				Items []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"items"`
			}

			if err := json.Unmarshal(bodyPlaylists, &playlistsResponse); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error unmarshalling playlists response:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling playlists response\" }\n")
				return
			}

			var removedFromPlaylist bool
			for _, playlist := range playlistsResponse.Items {
				playlistID := playlist.ID
				playlistURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)
				reqPlaylistTracks, err := http.NewRequest("GET", playlistURL, nil)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println("Error creating playlist tracks request:", err.Error())
					fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
					return
				}

				reqPlaylistTracks.Header.Set("Authorization", "Bearer "+spotifyToken)
				respPlaylistTracks, err := client.Do(reqPlaylistTracks)
				if err != nil {
					w.WriteHeader(http.StatusBadGateway)
					fmt.Println("Error fetching playlist tracks:", err.Error())
					fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
					return
				}
				defer respPlaylistTracks.Body.Close()

				bodyPlaylistTracks, err := io.ReadAll(respPlaylistTracks.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println("Error reading playlist tracks response body:", err.Error())
					fmt.Fprintf(w, "{ \"error\": \"Error reading playlist tracks response body\" }\n")
					return
				}

				var playlistTracksResponse struct {
					Items []struct {
						Track struct {
							ID string `json:"id"`
						} `json:"track"`
					} `json:"items"`
				}

				if err := json.Unmarshal(bodyPlaylistTracks, &playlistTracksResponse); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println("Error unmarshalling playlist tracks response:", err.Error())
					fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling playlist tracks response\" }\n")
					return
				}

				for _, track := range playlistTracksResponse.Items {
					if track.Track.ID == trackID {
						removeTrackURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)
						removeTrackPayload := fmt.Sprintf(`{"tracks":[{"uri":"spotify:track:%s"}]}`, trackID)
						reqRemoveTrack, err := http.NewRequest("DELETE", removeTrackURL, bytes.NewBuffer([]byte(removeTrackPayload)))
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							fmt.Println("Error creating remove track request:", err.Error())
							fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
							return
						}

						reqRemoveTrack.Header.Set("Authorization", "Bearer "+spotifyToken)
						respRemoveTrack, err := client.Do(reqRemoveTrack)
						if err != nil {
							w.WriteHeader(http.StatusBadGateway)
							fmt.Println("Error removing track from playlist:", err.Error())
							fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
							return
						}
						defer respRemoveTrack.Body.Close()

						if respRemoveTrack.StatusCode == http.StatusOK {
							fmt.Printf("Removed %s from playlist %s\n", trackName, playlist.Name)
							removedFromPlaylist = true
							break
						} else {
							fmt.Printf("Failed to remove track from playlist %s\n", playlist.Name)
						}
					}
				}
				if removedFromPlaylist {
					break
				}
			}

			if removedFromPlaylist {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "{ \"status\": \"Track removed from playlist\" }\n")
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "{ \"error\": \"Track not found in any playlist\" }\n")
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{ \"error\": \"No track is currently playing\" }\n")
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error fetching currently playing track. Status:", respSpotify.StatusCode)
		fmt.Fprintf(w, "{ \"error\": \"Failed to get currently playing track\" }\n")
	}
}

func addToPlaylistIfNotPresent(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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

	spotifyURL := "https://api.spotify.com/v1/me/player/currently-playing"
	reqSpotify, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error creating Spotify request:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqSpotify.Header.Set("Authorization", "Bearer "+spotifyToken)
	client := &http.Client{}
	respSpotify, err := client.Do(reqSpotify)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Println("Error fetching currently playing music:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer respSpotify.Body.Close()

	bodyResp, err := io.ReadAll(respSpotify.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error reading Spotify response body:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"Error reading Spotify response body\" }\n")
		return
	}
	fmt.Println("Spotify Response Body:", string(bodyResp))

	if respSpotify.StatusCode == http.StatusOK {
		var playbackResponse struct {
			Item struct {
				Name    string `json:"name"`
				ID      string `json:"id"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
			} `json:"item"`
		}

		if err := json.Unmarshal(bodyResp, &playbackResponse); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error unmarshalling playback response:", err.Error())
			fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling playback response\" }\n")
			return
		}

		if playbackResponse.Item.Name != "" {
			trackID := playbackResponse.Item.ID
			trackName := playbackResponse.Item.Name
			artistName := playbackResponse.Item.Artists[0].Name

			fmt.Printf("Currently playing: %s by %s\n", trackName, artistName)

			playlistsURL := "https://api.spotify.com/v1/me/playlists"
			reqPlaylists, err := http.NewRequest("GET", playlistsURL, nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error creating playlists request:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
				return
			}

			reqPlaylists.Header.Set("Authorization", "Bearer "+spotifyToken)
			respPlaylists, err := client.Do(reqPlaylists)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				fmt.Println("Error fetching playlists:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
				return
			}
			defer respPlaylists.Body.Close()

			bodyPlaylists, err := io.ReadAll(respPlaylists.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error reading playlists response body:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"Error reading playlists response body\" }\n")
				return
			}

			var playlistsResponse struct {
				Items []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"items"`
			}

			if err := json.Unmarshal(bodyPlaylists, &playlistsResponse); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error unmarshalling playlists response:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling playlists response\" }\n")
				return
			}

			var trackInAnyPlaylist bool
			for _, playlist := range playlistsResponse.Items {
				playlistID := playlist.ID
				playlistURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)
				reqPlaylistTracks, err := http.NewRequest("GET", playlistURL, nil)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println("Error creating playlist tracks request:", err.Error())
					fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
					return
				}

				reqPlaylistTracks.Header.Set("Authorization", "Bearer "+spotifyToken)
				respPlaylistTracks, err := client.Do(reqPlaylistTracks)
				if err != nil {
					w.WriteHeader(http.StatusBadGateway)
					fmt.Println("Error fetching playlist tracks:", err.Error())
					fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
					return
				}
				defer respPlaylistTracks.Body.Close()

				bodyPlaylistTracks, err := io.ReadAll(respPlaylistTracks.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println("Error reading playlist tracks response body:", err.Error())
					fmt.Fprintf(w, "{ \"error\": \"Error reading playlist tracks response body\" }\n")
					return
				}

				var playlistTracksResponse struct {
					Items []struct {
						Track struct {
							ID string `json:"id"`
						} `json:"track"`
					} `json:"items"`
				}

				if err := json.Unmarshal(bodyPlaylistTracks, &playlistTracksResponse); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println("Error unmarshalling playlist tracks response:", err.Error())
					fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling playlist tracks response\" }\n")
					return
				}

				for _, track := range playlistTracksResponse.Items {
					if track.Track.ID == trackID {
						trackInAnyPlaylist = true
						break
					}
				}

				if trackInAnyPlaylist {
					break
				}
			}

			if !trackInAnyPlaylist && len(playlistsResponse.Items) > 0 {
				firstPlaylistID := playlistsResponse.Items[0].ID
				addTrackURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", firstPlaylistID)
				addTrackPayload := fmt.Sprintf(`{"uris":["spotify:track:%s"]}`, trackID)
				reqAddTrack, err := http.NewRequest("POST", addTrackURL, bytes.NewBuffer([]byte(addTrackPayload)))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println("Error creating add track request:", err.Error())
					fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
					return
				}

				reqAddTrack.Header.Set("Authorization", "Bearer "+spotifyToken)
				respAddTrack, err := client.Do(reqAddTrack)
				if err != nil {
					w.WriteHeader(http.StatusBadGateway)
					fmt.Println("Error adding track to playlist:", err.Error())
					fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
					return
				}
				defer respAddTrack.Body.Close()

				if respAddTrack.StatusCode == http.StatusOK {
					fmt.Println("Successfully added track to the playlist")
					w.WriteHeader(http.StatusOK)
					fmt.Fprintf(w, "{ \"status\": \"Track added to the playlist\" }\n")
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println("Failed to add track to playlist. Status:", respAddTrack.StatusCode)
					fmt.Fprintf(w, "{ \"error\": \"Failed to add track to playlist\" }\n")
				}
			} else {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "{ \"status\": \"Track is already in a playlist or no playlists available\" }\n")
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{ \"error\": \"No track is currently playing\" }\n")
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error fetching currently playing track. Status:", respSpotify.StatusCode)
		fmt.Fprintf(w, "{ \"error\": \"Failed to get currently playing track\" }\n")
	}
}

func createPlaylist(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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
			Musique string `json:"name"`
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
	playlistName := requestBody.Spices.Musique
	fmt.Println("Extracted userID:", userID)
	fmt.Println("Playlist Name:", playlistName)

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

	spotifyURL := "https://api.spotify.com/v1/me/playlists"
	createPlaylistPayload := fmt.Sprintf(`{
        "name": "%s",
        "description": "A new playlist created via API",
        "public": false
    }`, playlistName)

	reqSpotify, err := http.NewRequest("POST", spotifyURL, bytes.NewBuffer([]byte(createPlaylistPayload)))
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
		fmt.Println("Error creating playlist:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer respSpotify.Body.Close()

	bodyResp, err := io.ReadAll(respSpotify.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error reading Spotify response body:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"Error reading Spotify response body\" }\n")
		return
	}
	fmt.Println("Spotify Response Body:", string(bodyResp))

	if respSpotify.StatusCode == http.StatusCreated {
		fmt.Println("Playlist successfully created!")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{ \"status\": \"Playlist successfully created!\" }\n")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed to create playlist. Status:", respSpotify.StatusCode)
		fmt.Fprintf(w, "{ \"error\": \"Failed to create playlist\" }\n")
	}
}

func clearPlaylist(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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
			PlaylistName string `json:"name"`
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
	playlistName := requestBody.Spices.PlaylistName
	fmt.Println("Extracted userID:", userID)
	fmt.Println("Playlist Name to clear:", playlistName)

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

	searchURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=playlist", playlistName)
	reqSearch, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error creating Spotify search request:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqSearch.Header.Set("Authorization", "Bearer "+spotifyToken)
	reqSearch.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	respSearch, err := client.Do(reqSearch)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Println("Error searching for playlist:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer respSearch.Body.Close()

	bodySearch, err := io.ReadAll(respSearch.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error reading Spotify search response body:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"Error reading Spotify search response body\" }\n")
		return
	}

	if respSearch.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed to search for playlist. Status:", respSearch.StatusCode)
		fmt.Fprintf(w, "{ \"error\": \"Failed to search for playlist\" }\n")
		return
	}

	var searchResult struct {
		Playlists struct {
			Items []struct {
				ID string `json:"id"`
			} `json:"items"`
		} `json:"playlists"`
	}

	if err := json.Unmarshal(bodySearch, &searchResult); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error unmarshalling search result:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling search result\" }\n")
		return
	}

	if len(searchResult.Playlists.Items) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("No playlist found with the name:", playlistName)
		fmt.Fprintf(w, "{ \"error\": \"No playlist found with the specified name\" }\n")
		return
	}

	playlistID := searchResult.Playlists.Items[0].ID
	fmt.Println("Found playlist ID:", playlistID)

	spotifyURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)

	reqSpotify, err := http.NewRequest("GET", spotifyURL, nil)
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
		fmt.Println("Error fetching playlist tracks:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer respSpotify.Body.Close()

	bodyResp, err := io.ReadAll(respSpotify.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error reading Spotify response body:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"Error reading Spotify response body\" }\n")
		return
	}
	fmt.Println("Spotify Response Body:", string(bodyResp))

	if respSpotify.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed to fetch playlist tracks. Status:", respSpotify.StatusCode)
		fmt.Fprintf(w, "{ \"error\": \"Failed to fetch playlist tracks\" }\n")
		return
	}

	var playlistTracks struct {
		Items []struct {
			Track struct {
				ID string `json:"id"`
			} `json:"track"`
		} `json:"items"`
	}

	if err := json.Unmarshal(bodyResp, &playlistTracks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error unmarshalling playlist tracks:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling playlist tracks\" }\n")
		return
	}

	if len(playlistTracks.Items) == 0 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{ \"status\": \"Playlist is already empty\" }\n")
		return
	}

	trackIDs := []string{}
	for _, item := range playlistTracks.Items {
		trackIDs = append(trackIDs, item.Track.ID)
	}

	removeTracksPayload := fmt.Sprintf(`{
        "tracks": [%s]
    }`, strings.Join(trackIDs, ","))

	reqSpotifyRemove, err := http.NewRequest("DELETE", spotifyURL, bytes.NewBuffer([]byte(removeTracksPayload)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error creating delete request:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqSpotifyRemove.Header.Set("Authorization", "Bearer "+spotifyToken)
	reqSpotifyRemove.Header.Set("Content-Type", "application/json")

	respSpotifyRemove, err := client.Do(reqSpotifyRemove)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Println("Error clearing playlist:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer respSpotifyRemove.Body.Close()

	bodyRespRemove, err := io.ReadAll(respSpotifyRemove.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error reading Spotify response body:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"Error reading Spotify response body\" }\n")
		return
	}
	fmt.Println("Spotify Response Body after removal:", string(bodyRespRemove))

	if respSpotifyRemove.StatusCode == http.StatusOK {
		fmt.Println("Playlist successfully cleared!")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{ \"status\": \"Playlist successfully cleared!\" }\n")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed to clear playlist. Status:", respSpotifyRemove.StatusCode)
		fmt.Fprintf(w, "{ \"error\": \"Failed to clear playlist\" }\n")
	}
}

func addToPlaylistByName(w http.ResponseWriter, req *http.Request, db *sql.DB) {
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
			Name string `json:"name"`
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
	playlistName := requestBody.Spices.Name
	fmt.Println("Extracted userID:", userID)
	fmt.Println("Playlist Name:", playlistName)

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

	spotifyURL := "https://api.spotify.com/v1/me/player/currently-playing"
	reqSpotify, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error creating Spotify request:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}

	reqSpotify.Header.Set("Authorization", "Bearer "+spotifyToken)
	client := &http.Client{}
	respSpotify, err := client.Do(reqSpotify)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Println("Error fetching currently playing music:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
		return
	}
	defer respSpotify.Body.Close()

	bodyResp, err := io.ReadAll(respSpotify.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error reading Spotify response body:", err.Error())
		fmt.Fprintf(w, "{ \"error\": \"Error reading Spotify response body\" }\n")
		return
	}
	fmt.Println("Spotify Response Body:", string(bodyResp))

	if respSpotify.StatusCode == http.StatusOK {
		var playbackResponse struct {
			Item struct {
				Name    string `json:"name"`
				ID      string `json:"id"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
			} `json:"item"`
		}

		if err := json.Unmarshal(bodyResp, &playbackResponse); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error unmarshalling playback response:", err.Error())
			fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling playback response\" }\n")
			return
		}

		if playbackResponse.Item.Name != "" {
			trackID := playbackResponse.Item.ID
			trackName := playbackResponse.Item.Name
			artistName := playbackResponse.Item.Artists[0].Name

			fmt.Printf("Currently playing: %s by %s\n", trackName, artistName)

			playlistsURL := "https://api.spotify.com/v1/me/playlists"
			reqPlaylists, err := http.NewRequest("GET", playlistsURL, nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error creating playlists request:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
				return
			}

			reqPlaylists.Header.Set("Authorization", "Bearer "+spotifyToken)
			respPlaylists, err := client.Do(reqPlaylists)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				fmt.Println("Error fetching playlists:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
				return
			}
			defer respPlaylists.Body.Close()

			bodyPlaylists, err := io.ReadAll(respPlaylists.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error reading playlists response body:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"Error reading playlists response body\" }\n")
				return
			}

			var playlistsResponse struct {
				Items []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"items"`
			}

			if err := json.Unmarshal(bodyPlaylists, &playlistsResponse); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error unmarshalling playlists response:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"Error unmarshalling playlists response\" }\n")
				return
			}

			var playlistID string
			for _, playlist := range playlistsResponse.Items {
				if playlist.Name == playlistName {
					playlistID = playlist.ID
					break
				}
			}

			if playlistID == "" {
				w.WriteHeader(http.StatusNotFound)
				fmt.Println("Playlist not found:", playlistName)
				fmt.Fprintf(w, "{ \"error\": \"Playlist not found\" }\n")
				return
			}

			addTrackURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)
			addTrackPayload := fmt.Sprintf(`{"uris":["spotify:track:%s"]}`, trackID)
			reqAddTrack, err := http.NewRequest("POST", addTrackURL, bytes.NewBuffer([]byte(addTrackPayload)))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error creating add track request:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
				return
			}

			reqAddTrack.Header.Set("Authorization", "Bearer "+spotifyToken)
			respAddTrack, err := client.Do(reqAddTrack)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				fmt.Println("Error adding track to playlist:", err.Error())
				fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
				return
			}
			defer respAddTrack.Body.Close()

			if respAddTrack.StatusCode == http.StatusOK {
				fmt.Println("Successfully added track to the playlist")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "{ \"status\": \"Track added to the playlist\" }\n")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Failed to add track to playlist. Status:", respAddTrack.StatusCode)
				fmt.Fprintf(w, "{ \"error\": \"Failed to add track to playlist\" }\n")
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{ \"error\": \"No track is currently playing\" }\n")
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error fetching currently playing track. Status:", respSpotify.StatusCode)
		fmt.Fprintf(w, "{ \"error\": \"Failed to get currently playing track\" }\n")
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
			Name: "playMusic",
			Type: "reaction",
			Desc: "Chose a music to play !",
			Spices: []InfoSpice{
				{
					Name:  "musique",
					Type:  "text",
					Title: "The message you want to send.",
				},
			},
		},
		{
			Name:   "pauseMusic",
			Type:   "reaction",
			Desc:   "Stop the current music !",
			Spices: []InfoSpice{},
		},
		{
			Name:   "resumeMusic",
			Type:   "reaction",
			Desc:   "Resume the current music !",
			Spices: []InfoSpice{},
		},
		{
			Name:   "checkDeviceConnection",
			Type:   "action",
			Desc:   "check if you have a current spotify running !",
			Spices: []InfoSpice{},
		},
		{
			Name:   "checkSongRunning",
			Type:   "action",
			Desc:   "check if you have a song whose running and got id, type and current song !",
			Spices: []InfoSpice{},
		},
		{
			Name:   "checkPodcastRunning",
			Type:   "action",
			Desc:   "check if you have a podcast whose running and got the podcast name and episode + follow him !",
			Spices: []InfoSpice{},
		},
		{
			Name:   "likeCurrentMusic",
			Type:   "reaction",
			Desc:   "like the current music you listening",
			Spices: []InfoSpice{},
		},
		{
			Name:   "unlikeCurrentMusic",
			Type:   "reaction",
			Desc:   "unlike the current music you listening",
			Spices: []InfoSpice{},
		},
		{
			Name:   "nextMusic",
			Type:   "reaction",
			Desc:   "next music",
			Spices: []InfoSpice{},
		},
		{
			Name:   "previousMusic",
			Type:   "reaction",
			Desc:   "previous music",
			Spices: []InfoSpice{},
		},
		{
			Name:   "removeFromPlaylistIfPresent",
			Type:   "reaction",
			Desc:   "remove the current music if she are in a playlist",
			Spices: []InfoSpice{},
		},
		{
			Name:   "addToPlaylistIfNotPresent",
			Type:   "reaction",
			Desc:   "add the current music if she are in a playlist",
			Spices: []InfoSpice{},
		},
		{
			Name: "createPlaylist",
			Type: "reaction",
			Desc: "creer une playlist",
			Spices: []InfoSpice{
				{
					Name:  "name",
					Type:  "text",
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
					Name:  "name",
					Type:  "text",
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
					Name:  "name",
					Type:  "text",
					Title: "The title of the playlist",
				},
			},
		},
	}
	var infos = Infos{
		Color:  "#1DB954",
		Image:  "/assets/spotify.png",
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
	fmt.Println("Spotify microservice container is running !")
	router := mux.NewRouter()
	godotenv.Load(".env")

	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/oauth", miniproxy(setOAUTHToken, db)).Methods("POST")
	router.HandleFunc("/playMusic", miniproxy(playMusic, db)).Methods("POST")
	router.HandleFunc("/pauseMusic", miniproxy(pauseMusic, db)).Methods("POST")
	router.HandleFunc("/resumeMusic", miniproxy(resumeMusic, db)).Methods("POST")
	router.HandleFunc("/checkDeviceConnection", miniproxy(checkDeviceConnection, db)).Methods("POST")
	router.HandleFunc("/checkSongRunning", miniproxy(checkSongRunning, db)).Methods("POST")
	router.HandleFunc("/checkPodcastRunning", miniproxy(checkPodcastRunning, db)).Methods("POST")
	router.HandleFunc("/likeCurrentMusic", miniproxy(likeCurrentMusic, db)).Methods("POST")
	router.HandleFunc("/unlikeCurrentMusic", miniproxy(unlikeCurrentMusic, db)).Methods("POST")
	router.HandleFunc("/nextMusic", miniproxy(nextMusic, db)).Methods("POST")
	router.HandleFunc("/previousMusic", miniproxy(previousMusic, db)).Methods("POST")
	router.HandleFunc("/removeFromPlaylistIfPresent", miniproxy(removeFromPlaylistIfPresent, db)).Methods("POST")
	router.HandleFunc("/addToPlaylistIfNotPresent", miniproxy(addToPlaylistIfNotPresent, db)).Methods("POST")
	router.HandleFunc("/createPlaylist", miniproxy(createPlaylist, db)).Methods("POST")
	router.HandleFunc("/clearPlaylist", miniproxy(clearPlaylist, db)).Methods("POST")
	router.HandleFunc("/addToPlaylistByName", miniproxy(addToPlaylistByName, db)).Methods("POST")
	router.HandleFunc("/user", getUserInfo).Methods("GET")
	router.HandleFunc("/", getRoutes).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}
