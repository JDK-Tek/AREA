package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	// "github.com/joho/godotenv"
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
			Area:    a,
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
	fmt.Println("bar")
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	obj, err := json.Marshal(message)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	fmt.Println("test")
	url := fmt.Sprintf("http://reverse-proxy:42002/service/%s/%s", service, name)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(obj))
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	fmt.Println("foo")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := http.Client{}
	rep, err := client.Do(req)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	defer rep.Body.Close()
	body, err := io.ReadAll(rep.Body)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	a.Reply(string(body), http.StatusOK)
}

func oauthGetter(a area.AreaRequest) {
	vars := mux.Vars(a.Request)
	service := vars["service"]
	redirect := a.Request.URL.Query().Get("redirect")
	url := fmt.Sprintf(
		"http://reverse-proxy:42002/service/%s/oauth?redirect=%s",
		service,
		url.QueryEscape(redirect),
	)
	fmt.Println(url)
	rep, err := http.Get(url)
	if err != nil {
		a.Error(err, rep.StatusCode)
		return
	}
	defer rep.Body.Close()
	if rep.StatusCode != http.StatusOK {
		a.Reply(map[string]string{
			"status": http.StatusText(rep.StatusCode),
		}, rep.StatusCode)
		return
	}
	data, err := io.ReadAll(rep.Body)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
	}
	str := html.UnescapeString(string(data))
	fmt.Println(str)
	a.Writter.WriteHeader(rep.StatusCode)
	fmt.Fprintln(a.Writter, str)
}

func oauthSetter(a area.AreaRequest) {
	vars := mux.Vars(a.Request)
	service := vars["service"]

	url := fmt.Sprintf(
		"http://reverse-proxy:42002/service/%s/oauth",
		service,
	)

	req, err := http.NewRequest("POST", url, a.Request.Body)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := http.Client{}
	rep, err := client.Do(req)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	defer rep.Body.Close()

	body, err := io.ReadAll(rep.Body)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	a.Reply(jsonResponse, http.StatusOK)
}

func codeCallback(a area.AreaRequest) {
	var code = a.Request.URL.Query().Get("code")

	if code == "" {
		a.ErrorStr("no code :(", http.StatusBadRequest)
		return
	}
	a.Reply(map[string]string{
		"message": "ok",
	}, http.StatusOK)
}

func getAllServices(a area.AreaRequest) {
	a.Reply(a.Area.Services, http.StatusOK)
}

func createTheAbout(a area.AreaRequest) {
	a.Area.About.Server.CurrentTime = time.Now().Unix()
	a.Reply(a.Area.About, http.StatusOK)
}

func getRoutes(a area.AreaRequest) {
	vars := mux.Vars(a.Request)
	service := vars["service"]
	url := fmt.Sprintf(
		"http://reverse-proxy:%s/service/%s/",
		os.Getenv("REVERSEPROXY_PORT"),
		service,
	)
	rep, err := http.Get(url)
	if err != nil {
		a.Error(err, rep.StatusCode)
		return
	}
	defer rep.Body.Close()
	body, err := io.ReadAll(rep.Body)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	if rep.StatusCode != 200 {
		a.ErrorStr(string(body), rep.StatusCode)
		return
	}
	a.Writter.WriteHeader(http.StatusOK)
	a.Writter.Write(body)
}

type Message struct {
	Message         string `json:"message"`
	Authentificated bool   `json:"authentificated"`
}

type MessageWithID struct {
	Message         string `json:"message"`
	Authentificated bool   `json:"authentificated"`
	ID              int    `json:"id"`
}

func doctor(a area.AreaRequest) {
	id, err := a.AssertToken()
	if err != nil {
		a.Reply(Message{Message: "i'm ok thanks", Authentificated: false}, http.StatusOK)
		return
	}
	a.Reply(MessageWithID{Message: "i'm ok thanks", Authentificated: true, ID: id}, http.StatusOK)
}

func main() {
	router := mux.NewRouter()
	var err error
	var dbPassword, dbUser, dbName, dbHost, dbPort string
	var servicePath string
	var connectStr string
	var portString string
	var a area.Area

	// if err = godotenv.Load("/usr/mount.d/.env"); err != nil {
	//     log.Fatal("no .env")
	// }
	if dbPassword = os.Getenv("DB_PASSWORD"); dbPassword == "" {
		log.Fatal("DB_PASSWORD not found")
	}
	if dbUser = os.Getenv("DB_USER"); dbUser == "" {
		log.Fatal("DB_USER not found")
	}
	if dbName = os.Getenv("DB_NAME"); dbName == "" {
		log.Fatal("DB_NAME not found")
	}
	if dbHost = os.Getenv("DB_HOST"); dbHost == "" {
		log.Fatal("DB_HOST not found")
	}
	if dbPort = os.Getenv("DB_PORT"); dbPort == "" {
		log.Fatal("DB_PORT not found")
	}
	if portString = os.Getenv("BACKEND_PORT"); portString == "" {
		log.Fatal("BACKEND_PORT not found")
	}
	if a.Key = os.Getenv("BACKEND_KEY"); a.Key == "" {
		log.Fatal("BACKEND_KEY not found")
	}
	if _, err = strconv.Atoi(portString); err != nil {
		log.Fatal("atoi:", err)
	}
	connectStr = fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
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
	if servicePath = os.Getenv("SERVICES_PATH"); servicePath == "" {
		log.Fatal("BACKEND_KEY not found")
	}
	err = a.ObserveServices(servicePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = a.SetupTheAbout()
	if err != nil {
		log.Fatal(err.Error())
	}
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	router.HandleFunc("/api/login", newProxy(&a, auth.DoSomeLogin)).Methods("POST")
	router.HandleFunc("/api/register", newProxy(&a, auth.DoSomeRegister)).Methods("POST")
	router.HandleFunc("/api/area", newProxy(&a, arearoute.NewArea)).Methods("POST")
	router.HandleFunc("/api/v2/services", newProxy(&a, service.GetServices)).Methods("GET")
	router.HandleFunc("/api/v2/service/{id}", newProxy(&a, service.GetServiceApplets)).Methods("GET")
	router.HandleFunc("/api/oauth/{service}", newProxy(&a, oauthGetter)).Methods("GET")
	router.HandleFunc("/api/oauth/{service}", newProxy(&a, oauthSetter)).Methods("POST")
	router.HandleFunc("/api/applets", newProxy(&a, applet.GetApplets)).Methods("GET")
	router.HandleFunc("/api/orchestrator", newProxy(&a, onUpdate)).Methods("PUT")
	router.HandleFunc("/api/services", newProxy(&a, getAllServices)).Methods("GET")
	router.HandleFunc("/api/services/{service}", newProxy(&a, getRoutes)).Methods("GET")
	router.HandleFunc("/api/doctor", newProxy(&a, doctor)).Methods("GET")
	router.HandleFunc("/api/change", newProxy(&a, auth.DoSomeChangePassword)).Methods("PUT")
	router.HandleFunc("/about.json", newProxy(&a, createTheAbout)).Methods("GET")

	fmt.Println("=> server listens on port ", portString)
	log.Fatal(http.ListenAndServe(":"+portString, corsMiddleware(router)))
}
