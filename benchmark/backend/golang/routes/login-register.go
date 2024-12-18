package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const mailRegexStr = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
const expiration = 60 * 30

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
    var special = regexp.MustCompile(`[@$\\/<>*+-_:?!#\^]`).MatchString(str)

    return lower && upper && number && special
}

func createAToken(email string, password string) (string, error) {
    var secret = []byte(password)
    var claims jwt.Claims
    var token *jwt.Token

    claims = jwt.MapClaims{
        "email": email,
        "exp": time.Now().Add(time.Second * expiration).Unix(),
    }
    token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secret)
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
    var tokenString string

    mail, password, err = getCredentials(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    tokenString, err = createAToken(mail, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Println("Email:", mail)
    fmt.Println("Password:", password)
    // TODO: the db stuff
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "{ \"token\": \"%s\" }\n", tokenString)
}

func DoSomeLogin(w http.ResponseWriter, req *http.Request) {
    var err error
    var mail, password string
    var tokenString string

    mail, password, err = getCredentials(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    tokenString, err = createAToken(mail, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Println("Email:", mail)
    fmt.Println("Password:", password)
    // TODO: the db stuff
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "{ \"token\": \"%s\" }\n", tokenString)
}
