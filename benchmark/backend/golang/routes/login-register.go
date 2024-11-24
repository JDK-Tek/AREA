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
const secretKey = "a_temp_secret_i_will_change_to_env_later_inshallah"

type RegisterRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type RegisterError struct {
    message string
}

func (e *RegisterError) Error() string {
    return e.message
}


func checkForPassword(str string) bool {
    var lower = regexp.MustCompile(`[a-z]`).MatchString(str)
    var upper = regexp.MustCompile(`[A-Z]`).MatchString(str)
    var number = regexp.MustCompile(`[0-9]`).MatchString(str)
    var special = regexp.MustCompile(`[@$\\/<>*+:?!#\^]`).MatchString(str)

    return lower && upper && number && special
}

func getCredentials(req *http.Request) (string, string, error) {
    var data RegisterRequest
	var decoder = json.NewDecoder(req.Body)
    var err = decoder.Decode(&data)
    var hashedBytes [32]byte
    var hashedString string
    var mailRe *regexp.Regexp
    
	if err != nil {
        return "", "", &RegisterError{message: "Invalid input"}
    }
    if data.Email == "" || data.Password == "" {
        return "", "", &RegisterError{message: "Required email and password"}
    }
    mailRe = regexp.MustCompile(mailRegexStr)
    if !mailRe.MatchString(data.Email) {
        return "", "", &RegisterError{message: "Invalid email pattern"}
    }
    if !checkForPassword(data.Password) {
        return "", "", &RegisterError{message: "Password too weak"}
    }
	hashedBytes = sha256.Sum256([]byte(data.Password))
    hashedString = hex.EncodeToString(hashedBytes[:])
    return data.Email, hashedString, nil
}

func DoSomeRegister(w http.ResponseWriter, req *http.Request) {
    var err error
    var mail, password string

    mail, password, err = getCredentials(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Println("Received registration request:")
    fmt.Println("Email:", mail)
    fmt.Println("Password:", password)
    // TODO: the db stuff
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "Registration successful")
}

func DoSomeLogin(w http.ResponseWriter, req *http.Request) {
    var err error
    var mail, password string

    mail, password, err = getCredentials(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Println("Received registration request:")
    fmt.Println("Email:", mail)
    fmt.Println("Password:", password)
    // TODO: Check the password in db
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "Registration successful")
}