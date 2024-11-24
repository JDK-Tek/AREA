package routes

import (
    "regexp"
	"fmt"
	"net/http"
	"crypto/sha256"
	"encoding/json"
    "encoding/hex"
)

const mailRegexStr = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"

type RegisterRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

func checkForPassword(str string) bool {
    var lower = regexp.MustCompile(`[a-z]`).MatchString(str)
    var upper = regexp.MustCompile(`[A-Z]`).MatchString(str)
    var number = regexp.MustCompile(`[0-9]`).MatchString(str)
    var special = regexp.MustCompile(`[@$\\/<>*+:?!#\^]`).MatchString(str)

    return lower && upper && number && special
}

func DoSomeRegister(w http.ResponseWriter, req *http.Request) {
    var data RegisterRequest
	var decoder = json.NewDecoder(req.Body)
    var err = decoder.Decode(&data)
    var hashedBytes [32]byte
    var hashedString string
    var mailRe *regexp.Regexp
    
	if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    if data.Email == "" || data.Password == "" {
        http.Error(w, "Email and password are required", http.StatusBadRequest)
        return
    }
    mailRe = regexp.MustCompile(mailRegexStr)
    if !mailRe.MatchString(data.Email) {
        http.Error(w, "invalid email", http.StatusBadRequest)
        return
    }
    if !checkForPassword(data.Password) {
        http.Error(w, "password need to have lowercase, uppercase, numbers and special", http.StatusBadRequest)
        return
    }
	hashedBytes = sha256.Sum256([]byte(data.Password))
    hashedString = hex.EncodeToString(hashedBytes[:])
    fmt.Println("Received registration request:")
    fmt.Println("Email:", data.Email)
    fmt.Println("Password:", hashedString)
    // TODO: the db stuff
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "Registration successful")
}