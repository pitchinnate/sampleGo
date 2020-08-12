package main

import (
  "flag"
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "os"
  "strings"
  "time"
  "goSample/response"
  "fmt"
  "github.com/joho/godotenv"
)

func main() {
  var isProduction bool
  flag.BoolVar(&isProduction, "prod", false, "App is running in production")
  flag.Parse()
  if !isProduction {
    loadConfig()
  }
  runWebServer()
}

func loadConfig() {
  log.Println("See if we can access an .env file directory")
  err := godotenv.Load(".env")
  if err != nil {
    log.Fatal(err)
  }
}

func runWebServer() {
  host, ok := os.LookupEnv("HOST")
  if !ok {
    log.Fatalf("missing required env var HOST")
  }
  port, ok := os.LookupEnv("PORT")
  if !ok {
    log.Fatalf("missing required env var PORT")
  }
  fmt.Printf("Running Web Server on  %s:%s\n", host, port)

  router := mux.NewRouter()
  router.StrictSlash(false)

  router.HandleFunc("/hello", getHome).Methods("GET")
  router.NotFoundHandler = http.HandlerFunc(pageNotFound)

  srv := buildServer(router, host, port)
  log.Fatal(srv.ListenAndServe())
}

func getHome(w http.ResponseWriter, r *http.Request) {
  message := response.TextResponse{"Hello world"}
  response.RespondWithJSON(w, 200, message)
}

func pageNotFound(w http.ResponseWriter, r *http.Request) {
  message := response.TextResponse{Message: "PAGE NOT FOUND"}
  response.RespondWithJSON(w, 404, message)
}

func buildServer(router *mux.Router, host string, port string) *http.Server {
  srv := &http.Server{
    Handler:      logTrim(router),
    Addr:         fmt.Sprintf("%s:%s", host, port),
    WriteTimeout: 15 * time.Second,
    ReadTimeout:  15 * time.Second,
  }
  return srv
}

func logTrim(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    corsSetup(&w, r)
    if r.Method == "OPTIONS" {
      return
    }
    if r.URL.Path != "/" {
      r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
    }
    fmt.Println(fmt.Sprintf("%s %s by %s", r.Method, r.RequestURI, r.RemoteAddr))
    next.ServeHTTP(w, r)
  })
}

func corsSetup(w *http.ResponseWriter, r *http.Request) {
  (*w).Header().Set("Access-Control-Allow-Origin", "*")
  (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Auth-Token, Accept-Language")
}
