package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	//"io"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const API_OAUTH_SPOTIFY = "https://accounts.spotify.com/api/token"
const API_USER_SPOTIFY = "https://api.spotify.com/v1/me"

func getOAUTHLink(w http.ResponseWriter, req *http.Request) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")

	if clientID == "" || redirectURI == "" {
		http.Error(w, "Client ID ou URI de redirection manquant", http.StatusInternalServerError)
		return
	}

	authURL := fmt.Sprintf(
		"https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=user-library-read",
		clientID,
		url.QueryEscape(redirectURI),
	)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, authURL)
}

func setOAUTHToken(w http.ResponseWriter, req *http.Request) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		http.Error(w, "Client ID, Client Secret ou URI de redirection manquants", http.StatusInternalServerError)
		return
	}

	var res struct {
		Code string `json:"code"`
	}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&res)
	if err != nil {
		http.Error(w, "Erreur lors du décodage de la requête : "+err.Error(), http.StatusInternalServerError)
		return
	}

	if res.Code == "" {
		http.Error(w, "Code d'authentification manquant", http.StatusBadRequest)
		return
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", res.Code)
	data.Set("redirect_uri", redirectURI)

	resp, err := http.PostForm("https://accounts.spotify.com/api/token", data)
	if err != nil {
		http.Error(w, "Erreur lors de l'échange du code contre un token : "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Erreur de l'API Spotify : "+resp.Status, http.StatusInternalServerError)
		return
	}

	var tokenResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		http.Error(w, "Erreur lors de la décodification de la réponse : "+err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, ok := tokenResp["access_token"].(string)
	if !ok {
		http.Error(w, "Erreur : accès token manquant dans la réponse", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"access_token\": \"%s\"}\n", accessToken)
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

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()

	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/callback", setOAUTHToken).Methods("POST")
	router.HandleFunc("/user", getUserInfo).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
