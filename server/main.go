package main

import (
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"./models"
	"./handlers"
)

import gHandlers "github.com/gorilla/handlers"

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable port=32768",
		"bg", "qwerty", "bgdb")
	err := models.InitDB(dbinfo)
	if err != nil {
		panic("Failed to connect database")
	}

	r := mux.NewRouter()
	r.HandleFunc("/ws", handlers.WSHandler)
	r.HandleFunc("/api/v1/rooms/", handlers.RoomViewSet)
	r.HandleFunc("/api/v1/rooms/{id:[0-9]+}/", handlers.RoomViewSet)
	r.HandleFunc("/api/v1/games/", handlers.GameViewSet)
	r.HandleFunc("/api/v1/games/{id:[0-9]+}/", handlers.GameViewSet)
	r.HandleFunc("/api/v1/games/{id:[0-9]+}/statistics/", handlers.StatsViewSet)
	r.HandleFunc("/api/v1/gamers/{id:[0-9]+}/", handlers.GamerViewSet)

	r.HandleFunc("/api/v1/accounts/signup/", handlers.SignUpHandler)
	r.HandleFunc("/api/v1/accounts/signin/", handlers.SignInHandler)
	r.HandleFunc("/api/v1/accounts/signout/", handlers.SignOutHandler)
	r.HandleFunc("/api/v1/accounts/current/", handlers.CurrentHandler)

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("../ui/dist/"))))

	fmt.Println("Listening on: ", "localhost:8080")

	// TODO: investigate CORS, ...
	methods := []string{"POST", "PUT", "GET", "OPTIONS", "DELETE"}
	allowMethods := gHandlers.AllowedMethods(methods)
	allowCreds := gHandlers.AllowCredentials()

	err = http.ListenAndServe("127.0.0.1:8000", gHandlers.LoggingHandler(os.Stdout, gHandlers.CORS(allowMethods, allowCreds)(r)))
	if err != nil {
		log.Fatal("Server error", err.Error())
	}
}
