package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"io"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const API_OAUTH_SPOTIFY = "https://accounts.spotify.com/api/token"
const API_USER_SPOTIFY = "https://api.spotify.com/v1/me"

// Fonction pour obtenir le lien d'authentification OAuth pour Spotify
func getOAUTHLink(w http.ResponseWriter, req *http.Request) {
	// Charger les variables d'environnement
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

// Fonction pour récupérer le token d'accès après que l'utilisateur se soit authentifié
func setOAUTHToken(w http.ResponseWriter, req *http.Request) {
	// Charger les variables d'environnement
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		http.Error(w, "Client ID, Client Secret ou URI de redirection manquants", http.StatusInternalServerError)
		return
	}

	// Décode la requête pour récupérer le code
	var res struct {
		Code string `json:"code"`
	}

	// Récupérer le code envoyé dans le body de la requête POST
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&res)
	if err != nil {
		http.Error(w, "Erreur lors du décodage de la requête : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Si le code est vide, retourner une erreur
	if res.Code == "" {
		http.Error(w, "Code d'authentification manquant", http.StatusBadRequest)
		return
	}

	// Préparer les données pour obtenir le token d'accès
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", res.Code)
	data.Set("redirect_uri", redirectURI)

	// Faire la requête POST vers l'API OAuth de Spotify pour obtenir le token
	resp, err := http.PostForm("https://accounts.spotify.com/api/token", data)
	if err != nil {
		http.Error(w, "Erreur lors de l'échange du code contre un token : "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Vérifier que la réponse est correcte
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Erreur de l'API Spotify : "+resp.Status, http.StatusInternalServerError)
		return
	}

	// Décoder la réponse de Spotify pour obtenir le token
	var tokenResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		http.Error(w, "Erreur lors de la décodification de la réponse : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extraire le token d'accès de la réponse
	accessToken, ok := tokenResp["access_token"].(string)
	if !ok {
		http.Error(w, "Erreur : accès token manquant dans la réponse", http.StatusInternalServerError)
		return
	}

	// Retourner le token d'accès dans la réponse
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"access_token\": \"%s\"}\n", accessToken)
}


// Fonction pour obtenir des informations sur l'utilisateur à partir de Spotify
func getUserInfo(w http.ResponseWriter, req *http.Request) {
	// Récupérer le token d'accès de l'utilisateur depuis l'en-tête Authorization
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Token d'accès manquant", http.StatusUnauthorized)
		return
	}

	// Faire une requête à l'API Spotify pour obtenir les informations de l'utilisateur
	reqUser, err := http.NewRequest("GET", API_USER_SPOTIFY, nil)
	if err != nil {
		http.Error(w, "Erreur lors de la création de la requête utilisateur", http.StatusInternalServerError)
		return
	}

	// Ajouter le token d'accès à l'en-tête de la requête
	reqUser.Header.Set("Authorization", "Bearer "+authHeader[7:])

	client := &http.Client{}
	resp, err := client.Do(reqUser)
	if err != nil {
		http.Error(w, "Erreur lors de l'appel API Spotify", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Décoder la réponse pour obtenir les informations de l'utilisateur
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Erreur lors de la décodification des informations utilisateur", http.StatusInternalServerError)
		return
	}

	// Retourner les informations de l'utilisateur au format JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userInfo)
}

func main() {
	// Charger les variables d'environnement depuis le fichier .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Configuration du routeur
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/oauth", getOAUTHLink).Methods("GET")
	router.HandleFunc("/callback", setOAUTHToken).Methods("POST")
	router.HandleFunc("/user", getUserInfo).Methods("GET")

	// Démarrer le serveur
	log.Fatal(http.ListenAndServe(":8080", router))
}
