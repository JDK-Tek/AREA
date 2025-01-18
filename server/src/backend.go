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
	"regexp"
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
	BridgeID    int               `json:"bridge"`
	Id          int               `json:"userid"`
	Ingredients map[string]string `json:"ingredients"`
}

type UserMessage struct {
	Spices json.RawMessage `json:"spices"`
	Id     int             `json:"userid"`
}

func applyIngredients(str string, ingredients map[string]string) string {
	var re = regexp.MustCompile(`\{([\w\.]+)\}`)

	result := re.ReplaceAllStringFunc(str, func(match string) string {
		key := match[1 : len(match)-1]
		if value, found := ingredients[key]; found {
			return value
		}
		return match
	})
	return result
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

	fmt.Println(ureq.Ingredients)

	// fill the ingreients
	jsonStr, err := url.QueryUnescape(string(message.Spices))
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	processedStr := applyIngredients(jsonStr, ureq.Ingredients)
	message.Spices = json.RawMessage(processedStr)

	// fill the message
	message.Id = ureq.Id
	obj, err := json.Marshal(message)
	if err != nil {
		a.Error(err, http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("http://reverse-proxy:%s/service/%s/%s", os.Getenv("REVERSEPROXY_PORT"), service, name)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(obj))
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
	a.Reply(string(body), http.StatusOK)
}

func oauthGetter(a area.AreaRequest) {
	vars := mux.Vars(a.Request)
	service := vars["service"]
	redirect := a.Request.URL.Query().Get("redirect")
	url := fmt.Sprintf(
		"http://reverse-proxy:%s/service/%s/oauth?redirect=%s",
		os.Getenv("REVERSEPROXY_PORT"),
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
		// "http://reverse-proxy:42002/service/%s/oauth?%s",
		"http://reverse-proxy:%s/service/%s/oauth",
		os.Getenv("REVERSEPROXY_PORT"),
		service,
	)

	req, err := http.NewRequest("POST", url, a.Request.Body)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", a.Request.Header.Get("Authorization"))

	client := http.Client{}
	rep, err := client.Do(req)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	defer rep.Body.Close()

	a.Writter.Header().Set("Content-Type", rep.Header.Get("Content-Type"))
	a.Writter.WriteHeader(rep.StatusCode)
	_, err = io.Copy(a.Writter, rep.Body)
	if err != nil {
		http.Error(a.Writter, err.Error(), http.StatusInternalServerError)
	}

	// body, err := io.ReadAll(rep.Body)
	// if err != nil {
	// 	a.Error(err, http.StatusBadGateway)
	// 	return
	// }

	// var jsonResponse map[string]interface{}
	// if err := json.Unmarshal(body, &jsonResponse); err != nil {
	// 	a.Error(err, http.StatusInternalServerError)
	// 	return
	// }
	// a.Reply(jsonResponse, http.StatusOK)

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

type MiniAbout struct {
	Color string `json:"color"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Oauth bool   `json:"oauth"`
}

func getAllServices(a area.AreaRequest) {
	var tmp MiniAbout
	services := []MiniAbout{}

	for _, v := range a.Area.About.Server.Services {
		tmp.Color = v.Color
		tmp.Image = v.Icon
		tmp.Name = v.Name
		tmp.Oauth = v.Oauth
		services = append(services, tmp)
	}
	a.Reply(services, http.StatusOK)
}

func createTheAbout(a area.AreaRequest) {
	a.Area.About.Server.CurrentTime = time.Now().Unix()
	a.Reply(a.Area.About, http.StatusOK)
}

func getServiceRoutes(a area.AreaRequest) {
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

type MessageWithIDAndOauths struct {
	Message         string   `json:"message"`
	Authentificated bool     `json:"authentificated"`
	ID              int      `json:"id"`
	Oauths          []string `json:"oauths"`
}

func doctor(a area.AreaRequest) {
	// cheeck for no token
	id, err := a.AssertToken()
	if err != nil {
		a.Reply(Message{Message: "i'm ok thanks", Authentificated: false}, http.StatusOK)
		return
	}

	// check for oauths
	rows, err := a.Area.Database.Query("select service from tokens where owner = $1", id)
	if err != nil {
		a.Reply(Message{Message: "i'm ill: " + err.Error(), Authentificated: true}, http.StatusOK)
		return
	}
	defer rows.Close()

	// check for no oauths
	var stuff string
	x := MessageWithIDAndOauths{Message: "i'm ok thanks", Authentificated: true, ID: id, Oauths: []string{}}
	for rows.Next() {
		if err := rows.Scan(&stuff); err != nil {
			a.Reply(Message{Message: "i'm ill: " + err.Error(), Authentificated: true}, http.StatusOK)
			return
		}
		fmt.Println(stuff)
		x.Oauths = append(x.Oauths, stuff)
	}
	if err := rows.Err(); err != nil {
		a.Reply(Message{Message: "i'm ill: " + err.Error(), Authentificated: true}, http.StatusOK)
		return
	}
	a.Reply(x, http.StatusOK)
}

func openWebhooks(a area.AreaRequest) {
	vars := mux.Vars(a.Request)
	service := vars["service"]
	toCall := fmt.Sprintf(
		"http://reverse-proxy:%s/service/%s/webhook",
		os.Getenv("REVERSEPROXY_PORT"),
		service,
	)

	// i copy all the querry parmas
	queryParams := a.Request.URL.Query()
	uri := fmt.Sprintf("%s?%s", toCall, queryParams.Encode())

	// i copy the body
	body, err := io.ReadAll(a.Request.Body)
	if err != nil {
		http.Error(a.Writter, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// i create a request
	req, err := http.NewRequest("POST", uri, bytes.NewReader(body))
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}

	// i copy the headers
	for k, vs := range a.Request.Header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	// snd the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// reecopy the headers
	for key, values := range resp.Header {
		for _, value := range values {
			a.Writter.Header().Add(key, value)
		}
	}

	// then i just like create thee response
	rep, err := io.ReadAll(resp.Body)
	if err != nil {
		a.Error(err, http.StatusBadGateway)
		return
	}
	a.Writter.WriteHeader(resp.StatusCode)
	a.Writter.Write(rep)
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
	router.HandleFunc("/api/services/{service}", newProxy(&a, getServiceRoutes)).Methods("GET")
	router.HandleFunc("/api/services/{service}/webhook", newProxy(&a, openWebhooks)).Methods("POST")
	router.HandleFunc("/api/doctor", newProxy(&a, doctor)).Methods("GET")
	router.HandleFunc("/api/change", newProxy(&a, auth.DoSomeChangePassword)).Methods("PUT")
	router.HandleFunc("/about.json", newProxy(&a, createTheAbout)).Methods("GET")
	router.PathPrefix("/assets").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	fmt.Println("=> server listens on port ", portString)
	time.Sleep(time.Second * 2)
	err = a.SetupTheAbout()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Fatal(http.ListenAndServe(":"+portString, corsMiddleware(router)))

}
