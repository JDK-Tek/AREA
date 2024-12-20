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
		"postgresql://%s:%s@database:5432/area_database?sslmode=disable",
		dbUser,
		dbPassword,
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
}

func masterThread(db *sql.DB) {
	var msg Message

	client := http.Client{}
	url := "http://backend:42000/api/orchestrator"
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
		for _, v := range bridges {
			msg.Bridge = v
			obj, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			req, err := http.NewRequest("PUT", url, bytes.NewBuffer(obj))
			fmt.Println("object is", string(obj))
			if err != nil {
				continue
			}
			rep, err := client.Do(req)
			if err != nil {
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
	log.Fatal(http.ListenAndServe(":80", router))
}