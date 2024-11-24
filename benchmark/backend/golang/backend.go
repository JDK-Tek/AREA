package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"backend-benchmark/routes"
)

const PORT int = 1234

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

func main() {
	var portString string = strconv.Itoa(PORT)
	var router *mux.Router = mux.NewRouter()

	fmt.Println("=> server listens on port ", PORT)
	router.HandleFunc("/hello/{name}", doSomeHello).Methods("GET")
	router.HandleFunc("/api/register", routes.DoSomeRegister).Methods("POST")
	router.HandleFunc("/api/login", routes.DoSomeLogin).Methods("POST")
	log.Fatal(http.ListenAndServe(":" + portString, router))
}