package auth

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

    "area-backend/area"
)

const mailRegexStr = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"

type RegisterRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type ChangePassword struct {
    Email string `json:"email"`
    Password string `json:"password"`
    NewPassword string `json:"new_password"`
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

func DoSomeRegister(a area.AreaRequest) {
    var err error
    var mail, password string
    var tokenString string
    var userid = -1

    mail, password, err = getCredentials(a.Request)
    if err != nil {
        a.Error(err, http.StatusBadRequest)
        return
    }
    fmt.Println("Email:", mail)
    fmt.Println("Password:", password)
    err = a.Area.Database.QueryRow("SELECT id FROM users WHERE email = $1", mail).Scan(&userid)
    if err != nil && err != sql.ErrNoRows {
        a.Error(err, http.StatusInternalServerError)
        return
    }
    if userid != -1 {
        a.ErrorStr("user already exists", http.StatusBadRequest)
        return
    }
    err = a.Area.Database.
        QueryRow("INSERT INTO users (email, password) VALUES ($1, $2) returning id", mail, password).
        Scan(&userid)
    if err != nil {
        a.Error(err, http.StatusInternalServerError)
        return
    }
    fmt.Println(userid)
    tokenString, err = a.Area.NewToken(userid)
    if err != nil {
        a.Error(err, http.StatusInternalServerError)
        return
    }
    a.Reply(map[string]any{
        "token": tokenString,
    }, http.StatusOK)
}

func DoSomeLogin(a area.AreaRequest) {
    var err error
    var mail, password, realPassword string
    var tokenString string
    var userid int

    mail, password, err = getCredentials(a.Request)
    if err != nil {
        a.Error(err, http.StatusBadRequest)
        return
    }
    fmt.Println("Email:", mail)
    fmt.Println("Password:", password)
    err = a.Area.Database.QueryRow("SELECT password, id FROM users WHERE email = $1", mail).Scan(&realPassword, &userid)
    if err != nil {
        a.ErrorStr("invalid email/password", http.StatusBadRequest)
        return
    }
    if realPassword != password {
        a.ErrorStr("invalid email/password", http.StatusBadRequest)
        return
    }
    tokenString, err = a.Area.NewToken(userid)
    if err != nil {
        a.Error(err, http.StatusInternalServerError)
        return
    }
    a.Reply(map[string]any{
        "token": tokenString,
    }, http.StatusOK)
}

func DoSomeChangePassword(a area.AreaRequest) {
    var realid int

    mail, password, err := getCredentials(a.Request)
    if err != nil {
        a.Error(err, http.StatusBadRequest)
        return
    }
    id, err := a.AssertToken()
    if err != nil {
        a.Error(err, http.StatusBadRequest)
        return
    }
    err = a.Area.Database.QueryRow(
        "select id from users where id = $1 and email = $2", id, mail).Scan(&realid)
    fmt.Println("hello")
    if err != nil {
        a.Error(err, http.StatusBadRequest)
        return
    }
    _, err = a.Area.Database.Exec("update users set password = $1 where id = $2", password, id)
    if err != nil {
        a.Error(err, http.StatusBadRequest)
        return
    }
    tokenString, err := a.Area.NewToken(id)
    if err != nil {
        a.Error(err, http.StatusInternalServerError)
        return
    }
    a.Reply(map[string]any{
        "token": tokenString,
    }, http.StatusOK)
}
