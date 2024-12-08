package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"area-backend/area"
	"area-backend/routes/applet"
	"area-backend/routes/arearoute"
	"area-backend/routes/auth"
	"area-backend/routes/service"
)

func newProxy(a *area.Area, f func(area.AreaRequest)) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        f(area.AreaRequest{
            Area: a,
            Writter: w,
            Request: r,
        })
    }
}

type UpdateRequest struct {
	BridgeID int `json:"bridge"`
}

type UserMessage struct {
	Spices json.RawMessage `json:"spices"`
}

func onUpdate(a area.AreaRequest) {
	var message UserMessage
	var ureq UpdateRequest
	var reactid int
	var service, name string

	err := json.NewDecoder(a.Request.Body).Decode(&ureq)
	if err != nil {
		a.Error(err, http.StatusBadRequest)
		return
	}
	err = a.Area.Database.QueryRow("select reaction from bridge where id = $1", ureq.BridgeID).Scan(&reactid)
	if err != nil {
		a.Error(err, http.StatusBadRequest)
		return
	}
	err = a.Area.Database.
		QueryRow("select service, name, spices from reactions where id = $1", reactid).
		Scan(&service, &name, &message.Spices)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	obj, err := json.Marshal(message)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	fmt.Println(string(obj))
	url := fmt.Sprintf("http://reverse-proxy:42002/service/%s/%s", service, name)
	fmt.Println(service, name)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(obj))
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	fmt.Println("world")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := http.Client{}
	rep, err := client.Do(req)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	defer rep.Body.Close()
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	a.Reply(string(body), http.StatusOK)
}

type UpdateRequest struct {
	BridgeID int `json:"bridge"`
}

type UserMessage struct {
	Spices json.RawMessage `json:"spices"`
}

func onUpdate(a area.AreaRequest) {
	var message UserMessage
	var ureq UpdateRequest
	var reactid int

	url := "http://reverse-proxy:42002/service/discord/send"
	err := json.NewDecoder(a.Request.Body).Decode(&ureq)
	if err != nil {
		a.Error(err, http.StatusBadRequest)
		return
	}
	err = a.Area.Database.QueryRow("select reaction from bridge where id = $1", ureq.BridgeID).Scan(&reactid)
	if err != nil {
		a.Error(err, http.StatusBadRequest)
		return
	}
	err = a.Area.Database.QueryRow("select spices from reactions where id = $1", reactid).Scan(&message.Spices)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	obj, err := json.Marshal(message)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	fmt.Println(string(message.Spices))
	// spices = []byte(`{ "spices": { "channel": 1304486924299272194, "message": "bite" } }`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(obj))
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	fmt.Println("world")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := http.Client{}
	rep, err := client.Do(req)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	defer rep.Body.Close()
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	fmt.Println("Response body:", string(body))
}

func main() {
    router := mux.NewRouter()
    var err error
    var dbPassword, dbUser string
    var connectStr string
    var portString string
    var a area.Area

    if err = godotenv.Load("/usr/mount.d/.env"); err != nil {
        log.Fatal("no .env")
    }
    if dbPassword = os.Getenv("DB_PASSWORD"); dbPassword == "" {
        log.Fatal("DB_PASSWORD not found")
    }
    if dbUser = os.Getenv("DB_USER"); dbUser == "" {
        log.Fatal("DB_USER not found")
    }
    if portString = os.Getenv("BACKEND_PORT"); portString == "" {
        log.Fatal("BACKEND_PORT not found")
    }
    if a.Key = os.Getenv("BACKEND_KEY"); portString == "" {
        log.Fatal("BACKEND_KEY not found")
    }
    if _, err = strconv.Atoi(portString); err != nil {
        log.Fatal("atoi:", err)
    }
    connectStr = fmt.Sprintf(
        "postgresql://%s:%s@database:5432/area_database?sslmode=disable",
        dbUser,
        dbPassword,
    )
    if a.Database, err = sql.Open("postgres", connectStr); err != nil {
        log.Fatal(err)
    }
    defer a.Database.Close()
    err = a.Database.Ping()
    for err != nil {
        fmt.Println("ping:", err)
        time.Sleep(time.Second)
        err = a.Database.Ping()
    }

    corsMiddleware := handlers.CORS(
        handlers.AllowedOrigins([]string{"*"}),
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
    )
	router.HandleFunc("/api/login", newProxy(&a, auth.DoSomeLogin)).Methods("POST")
	router.HandleFunc("/api/register", newProxy(&a, auth.DoSomeRegister)).Methods("POST")
	router.HandleFunc("/api/area", newProxy(&a, arearoute.NewArea)).Methods("POST")
	router.HandleFunc("/api/services", newProxy(&a, service.GetServices)).Methods("GET")
	router.HandleFunc("/api/service/{id}", newProxy(&a, service.GetServiceApplets)).Methods("GET")
	router.HandleFunc("/api/applets", newProxy(&a, applet.GetApplets)).Methods("GET")
	router.HandleFunc("/api/orchestrator", newProxy(&a, onUpdate)).Methods("PUT")
    fmt.Println("=> server listens on port ", portString)
    log.Fatal(http.ListenAndServe(":"+portString, corsMiddleware(router)))
}
