package routes

import (
	"crypto/sha256"
	"database/sql"
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

func DoSomeRegister(w http.ResponseWriter, req *http.Request, db *sql.DB) {
    var err error
    var mail, password string
    var tokenString string
    var userid = -1

    mail, password, err = getCredentials(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Println("Email:", mail)
    fmt.Println("Password:", password)
    err = db.QueryRow("SELECT id FROM users WHERE email = $1", mail).Scan(&userid)
    if err != nil && err != sql.ErrNoRows {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if userid != -1 {
        http.Error(w, "user already exists", http.StatusBadRequest)
        return
    }
    tokenString, err = createAToken(mail, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    _, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", mail, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "{ \"token\": \"%s\" }\n", tokenString)
}

func DoSomeLogin(w http.ResponseWriter, req *http.Request, db *sql.DB) {
    var err error
    var mail, password, realPassword string
    var tokenString string

    mail, password, err = getCredentials(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    fmt.Println("Email:", mail)
    fmt.Println("Password:", password)
    err = db.QueryRow("SELECT password FROM users WHERE email = $1", mail).Scan(&realPassword)
    if err != nil {
        http.Error(w, "invalid user/password", http.StatusBadRequest)
        return
    }
    if realPassword != password {
        http.Error(w, "invalid user/password", http.StatusBadRequest)
        return
    }
    tokenString, err = createAToken(mail, password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "{ \"token\": \"%s\" }\n", tokenString)
}