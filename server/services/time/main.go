package main

import (
	"database/sql"
	"errors"
	"time"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	// "net/url"
	"os"
	// "strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const API = "https://tools.aimylogic.com/api/now?tz=Europe/Paris"

type Spices struct {
	HowMuch int `json:"howmuch"`
	Unit string `json:"unit"`
}

type Content struct {
	Spices Spices `json:"spices"`
	BridgeId int `json:"bridge"`
}

type Response struct {
	Timestamp int64 `json:"timestamp"`
}

func makeAnError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", err.Error())
}

func getTimeNow() (n int64, err error) {
	var r Response

	rep, err := http.Get(API)
	if err != nil {
		return
	}
	defer rep.Body.Close()
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return
	}
	n = r.Timestamp
	return
}

func getTimeAndContent(w http.ResponseWriter, req *http.Request) (c Content, n int64, err error) {
	n, err = getTimeNow()
	if err != nil {
		makeAnError(w, err, http.StatusBadGateway)
		return
	}
	c.BridgeId = -1
	c.Spices.HowMuch = -1
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&c)
	if err != nil {
		makeAnError(w, err, http.StatusBadRequest)
		return
	}
	if c.BridgeId == -1 || c.Spices.HowMuch == -1 {
		err = errors.New("invalid parsing")
		makeAnError(w, err, http.StatusBadRequest)
		return
	}
	return
}

func spices2Seconds(spices Spices) int64 {
	var association = map[string]int{
		"weeks": 7 * 24 * 3600,
		"days": 24 * 3600,
		"hours": 3600,
		"minutes": 60,
		"seconds": 1,
	}
	return int64(spices.HowMuch * association[spices.Unit])
}

func timeIn(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var maybeid int

	content, ms, err := getTimeAndContent(w, req)
	if err != nil {
		return
	}
	err = db.QueryRow("select id from micro_time where bridgeid = $1", content.BridgeId).Scan(&maybeid)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{ \"error\": \"it already exists\" }\n")
		return
	}
	secs := ms / 1000
	nsecs := (ms % 1000) * 1e6
	secs += spices2Seconds(content.Spices)
	timestamp := time.Unix(secs, nsecs)
	_, err = db.Exec("insert into micro_time (bridgeid, triggers) values ($1, $2)", content.BridgeId, timestamp)
	if err != nil {
		makeAnError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello %v", timestamp)
}

func connectToDatabase() (*sql.DB, error) {
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD not found")
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER not found")
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal("DB_HOST not found")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME not found")
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Fatal("DB_PORT not found")
	}
	// connectStr := fmt.Sprintf(
	// 	"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
	// 	dbUser,
	// 	dbPassword,
	// 	dbHost,
	// 	dbPort,
	// 	dbName,
	// )
	connectStr := fmt.Sprintf(
		"postgresql://%s:%s@database:%s/area_database?sslmode=disable",
		dbUser,
		dbPassword,
		dbPort,
	)
	return sql.Open("postgres", connectStr)
}

func miniProxy(f func(http.ResponseWriter, *http.Request, *sql.DB), c *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(a http.ResponseWriter, b *http.Request) {
		f(a, b, c)
	}
}

type Message struct {
	Bridge int `json:"bridge"`
	Ingredients map[string]string `json:"ingredients"`
}

func masterThread(db *sql.DB) {
	var msg Message

	msg.Ingredients = make(map[string]string)
	client := http.Client{}
	backendPort := os.Getenv("BACKEND_PORT")
	if backendPort == "" {
		log.Fatal("BACKEND_PORT not found")
	}
	url := fmt.Sprintf("http://backend:%s/api/orchestrator", backendPort)
	for {
		var bridges []int
		var n int

		time.Sleep(time.Second)
		ms, err := getTimeNow()
		if err != nil {
			fmt.Println("(x_x) <( here is what happend:", err.Error(), ")")
			continue
		}
		secs := ms / 1000
		nsecs := (ms % 1000) * 1e6
		timestamp := time.Unix(secs, nsecs)
		querry := "select bridgeid from micro_time where triggers < $1"
		rows, err := db.Query(querry, timestamp)
		if err != nil {
			fmt.Println("(x_x) <( here is what happend:", err.Error(), ")")
			continue
		}
		for rows.Next() {
			if err := rows.Scan(&n); err != nil {
				rows.Close()
				continue
			}
			bridges = append(bridges, n)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			continue
		}
		formatTime := timestamp.Format("2006-01-02 15:04:05")
		for _, v := range bridges {
			msg.Bridge = v
			msg.Ingredients["time"] = formatTime
			obj, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			req, err := http.NewRequest("PUT", url, bytes.NewBuffer(obj))
			fmt.Println("object is", string(obj))
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			rep, err := client.Do(req)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			rep.Body.Close()
			fmt.Println("finish triggering", v)
		}
		_, err = db.Exec("delete from micro_time where triggers < $1", timestamp)
		if err != nil {
			fmt.Println("(x_x) <( here is what happend:", err.Error(), ")")
			continue
		}
	}
}

type InfoSpice struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Title string `json:"title"`
	Extra []string `json:"extra"`
}

type InfoRoute struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Desc string `json:"description"`
	Spices []InfoSpice `json:"spices"`
}

type Infos struct {
	Color string `json:"color"`
	Image string `json:"image"`
	Routes []InfoRoute `json:"areas"`
}

func getRoutes(w http.ResponseWriter, req *http.Request) {
	var list = []InfoRoute{
		InfoRoute{
			Name: "in",
			Type: "action",
			Desc: "Triggers in some amount of time.",
			Spices: []InfoSpice{
				{
					Name: "howmuch",
					Type: "number",
					Title: "How much time to wait.",
				},
				{
					Name: "unit",
					Type: "dropdown",
					Title: "The unit to wait.",
					Extra: []string{
						"weeks",
						"days",
						"hours",
						"minutes",
						"seconds",
					},
				},
			},
		},
	}
	var infos = Infos{
		Color: "#ffffff",
		Image: "https://img.icons8.com/ios/452/timer.png",
		Routes: list,
	}
	var data []byte
	var err error

	data, err = json.Marshal(infos)
	if err != nil {
		http.Error(w, `{ "error": "marshal" }`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, string(data))
}

func main() {
	db, err := connectToDatabase()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(84)
	}
	go masterThread(db)
	fmt.Println("time microservice container is running !")
	router := mux.NewRouter()
	router.HandleFunc("/in", miniProxy(timeIn, db)).Methods("POST")
	router.HandleFunc("/", getRoutes).Methods("GET")
	log.Fatal(http.ListenAndServe(":80", router))
}