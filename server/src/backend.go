package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
	"os"
	"encoding/json"
	"database/sql"

	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"area-backend/routes"
	"area-backend/area"
)

func calculate(t0 []float64, t1 []float64, it int) (success bool, res float64, ms float64) {
	var start = time.Now()
	var vel = make([]float64, 3)
	var k float64
	var x uint64 = 0

	success = true
	for n := 0; n < it; n++ {
		x += uint64(n)
		if (t1[2] == t0[2]) || (t0[2] > 0 && t1[2] < 0) || (t0[2] < 0 && t1[2] > 0) {
			success = false
			continue
		}
		vel[0] = t1[0] - t0[0]
		vel[1] = t1[1] - t0[1]
		vel[2] = t1[2] - t0[2]
		k = math.Sqrt(vel[0] * vel[0] + vel[1] * vel[1] + vel[2] * vel[2])
		res = math.Asin(vel[2] / k)
		res = math.Floor(res * -180 / math.Pi * 100) / 100
	}
	fmt.Println(x)
	ms = float64(time.Since(start).Seconds() * 1000)
	return
}

func doSomeHello(w http.ResponseWriter, req *http.Request) {
	var vars = mux.Vars(req)
    var name = vars["name"]
	var arrays [2][]float64
	var err error

	err = json.Unmarshal([]byte(req.URL.Query().Get("t0")), &arrays[0])
	if err != nil {
		http.Error(w, "bad t0 parameter", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal([]byte(req.URL.Query().Get("t1")), &arrays[1])
	if err != nil {
		http.Error(w, "bad t1 parameter", http.StatusBadRequest)
		return
	}
	success, res, ms := calculate(arrays[0], arrays[1], 1_000_000)
	if !success {
		fmt.Fprintf(w, "Hello %s! your ball wont reach, computed in %.2fms", name, ms)
	} else {
		fmt.Fprintf(w, "Hello %s! your incidence angle is %.2f, computed in %.2fms", name, res, ms)
	}
}

func newProxy(a *area.Area, f func(area.AreaRequest)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(area.AreaRequest{
			Area: a,
			Writter: w,
			Request: r,
		})
	}
}

func main() {
	var router = mux.NewRouter()
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
	if err = a.Database.Ping(); err != nil {
		log.Fatal("ping:", err)
	}
	fmt.Println("=> server listens on port ", portString)
	router.HandleFunc("/hello/{name}", doSomeHello).Methods("GET")
	router.HandleFunc("/api/login", newProxy(&a, routes.DoSomeLogin)).Methods("POST")
	router.HandleFunc("/api/register", newProxy(&a, routes.DoSomeRegister)).Methods("POST")
	log.Fatal(http.ListenAndServe(":" + portString, router))
}