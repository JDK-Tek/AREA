package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"os"
	"database/sql"

	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"

	"area-backend/routes/auth"
	"area-backend/routes/arearoute"
	"area-backend/area"
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

func main() {
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
        handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Remplacez par l'origine de votre frontend
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
    )
    
    router := mux.NewRouter() // Assurez-vous que le routeur est initialisé avant d'appliquer le middleware
    
    router.HandleFunc("/api/login", newProxy(&a, auth.DoSomeLogin)).Methods("POST")
    router.HandleFunc("/api/register", newProxy(&a, auth.DoSomeRegister)).Methods("POST")
    router.HandleFunc("/api/area", newProxy(&a, arearoute.NewArea)).Methods("POST")
    
    fmt.Println("=> server listens on port ", portString)
    log.Fatal(http.ListenAndServe(":"+portString, corsMiddleware(router)))
}