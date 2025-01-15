package area

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const expiration = 60 * 30

// the about structure (for about.json)

type AboutSpice struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Title string `json:"title"`
	Extra []string `json:"extra"`
}

type AboutSomething struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Spices []AboutSpice `json:"spices"`
}

type AboutClient struct {
	Host string `json:"host"`
}

type AboutSevice struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Color string `json:"color"`
	Actions []AboutSomething `json:"actions"`
	Reactions []AboutSomething `json:"reactions"`
}

type AboutServer struct {
	CurrentTime int64 `json:"current_time"`
	Services []AboutSevice `json:"services"`
}

type About struct {
	Client AboutClient `json:"client"`
	Server AboutServer `json:"server"`
}

// the main area structuree

type Area struct {
	Database *sql.DB
	Key string
	Services []string
	About About
}

// for the informations i get from the services

type InfoRoute struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Desc string `json:"description"`
	Spices []AboutSpice `json:"spices"`
}

type Infos struct {
	Color string `json:"color"`
	Image string `json:"image"`
	Routes []InfoRoute `json:"areas"`
}

func (it *Area) ObserveServices(where string) error {
	entries, err := os.ReadDir(where)
	it.Services = make([]string, len(entries))
    if err != nil {
        return err
    }
    for n, e := range entries {
        it.Services[n] = e.Name()
    }
	return nil
}

func (it *Area) SetupTheAbout() error {
	var revproxy = os.Getenv("REVERSEPROXY_PORT")
	var infos Infos
	var tmpService AboutSevice
	var something AboutSomething

	it.About = About{
		Client: AboutClient{
			Host: os.Getenv("FRONTEND"),
		},
		Server: AboutServer{
			CurrentTime: time.Now().Unix(),
			Services: []AboutSevice{},
		},
	}
	if len(it.Services) == 0 {
		return nil
	}
	// fmt.Println(fmt.Sprintf("http://reverse-proxy:%s/service/%s/", revproxy, "coucou"))
	for _, service := range it.Services {
		url := fmt.Sprintf("http://reverse-proxy:%s/service/%s/", revproxy, service)
		rep, err := http.Get(url)
		if err != nil {
			fmt.Printf("%s failed: %s\n", service, err.Error())
			continue
		}
		defer rep.Body.Close()
		body, err := io.ReadAll(rep.Body)
		if err != nil {
			fmt.Printf("reeadall of %s failed: %s\n", service, err.Error())
			continue
		}
		// fmt.Printf("%s: (%s)\n", url, body)
		err = json.Unmarshal(body, &infos)
		if err != nil {
			fmt.Printf("decoding %s failed: %s body is %s\n", service, err.Error(), rep.Body)
			continue
		}
		tmpService.Actions = []AboutSomething{}
		tmpService.Reactions = []AboutSomething{}
		for _, v := range infos.Routes {
			something.Description = v.Desc
			something.Name = v.Name
			something.Spices = make([]AboutSpice, len(v.Spices))
			copy(something.Spices, v.Spices)
			if v.Type == "action" {
				tmpService.Actions = append(tmpService.Actions, something)
			} else {
				tmpService.Reactions = append(tmpService.Reactions, something)
			}
		}
		tmpService.Name = service
		tmpService.Icon = infos.Image
		tmpService.Color = infos.Color
		it.About.Server.Services = append(it.About.Server.Services, tmpService)
	}
	return nil
}

func (it *Area) NewToken(id int) (string, error) {
    var secretBytes = []byte(it.Key)
    var claims jwt.Claims
    var token *jwt.Token

    claims = jwt.MapClaims{
        "id": id,
        "exp": time.Now().Add(time.Second * expiration).Unix(),
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
		var id = claims["id"].(float64)
		return int(id), nil
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
