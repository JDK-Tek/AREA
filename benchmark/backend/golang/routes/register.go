package routes

import (
	"fmt"
	"net/http"
	"crypto/sha256"
	"encoding/json"
    "encoding/hex"
)

type RegisterRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

func DoSomeRegister(w http.ResponseWriter, req *http.Request) {
    var data RegisterRequest
	var decoder = json.NewDecoder(req.Body)
    var err = decoder.Decode(&data)
    var hashedBytes [32]byte
    var hashedString string
    
	if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    if data.Email == "" || data.Password == "" {
        http.Error(w, "Email and password are required", http.StatusBadRequest)
        return
    }
	hashedBytes = sha256.Sum256([]byte(data.Password))
    hashedString = hex.EncodeToString(hashedBytes[:])
    fmt.Println("Received registration request:")
    fmt.Println("Email:", data.Email)
    fmt.Println("Password:", hashedString)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Registration successful"))
}