package area

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const expiration = 60 * 30

type Area struct {
	Database *sql.DB
	Key string
}

func (it *Area) NewToken(id int) (string, error) {
    var secretBytes = []byte(it.Key)
    var claims jwt.Claims
    var token *jwt.Token

    claims = jwt.MapClaims{
        "id": id,
        "exp": time.Now().Add(time.Second * expiration).Unix(),
		"tokenid": -1,
    }
    token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretBytes)
}

func (it *Area) Token2Email(str string) (int, error) {
	var token, err = jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("bad method")
		}
		return []byte(it.Key), nil
	})
	var ok bool
	var claims jwt.MapClaims

	if err != nil {
		return -1, err
	}
	if claims, ok = token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["id"].(int)
		return id, nil
	}
	return -1, fmt.Errorf("invalid token or expired")
}

type AreaRequest struct {
	Area *Area
	Writter http.ResponseWriter
	Request *http.Request
}

func (it *AreaRequest) Error(err error, code int) {
	it.ErrorStr(err.Error(), code)
}

func (it *AreaRequest) ErrorStr(str string, code int) {
	var data []byte
	var err error

	data, err = json.Marshal(map[string]string{
		"error": str,
	})
	if err != nil {
		http.Error(it.Writter, "marshal failed", code)
		return
	}
	http.Error(it.Writter, string(data), code)
}

func (it *AreaRequest) Reply(object any, code int) {
	var data []byte
	var err error

	it.Writter.WriteHeader(http.StatusOK)
	data, err = json.Marshal(object)
	if err != nil {
		it.Error(err, http.StatusInternalServerError)
		return
	}
    fmt.Fprintln(it.Writter, string(data))
}

func (it *AreaRequest) AssertToken() (int, error) {
	var str = it.Request.Header.Get("Authorization")

	if str == "" || !strings.HasPrefix(str, "Bearer ") {
		return -1, fmt.Errorf("no authorization/bad init (no bearer maybe ?)")
	}
	str = strings.TrimPrefix(str, "Bearer ")
	return it.Area.Token2Email(str)
}
