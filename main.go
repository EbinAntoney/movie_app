package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"movies-app/api"
	"movies-app/helpers"
	"movies-app/middleware"
	"net/http"
)

func main() {
	fmt.Println("starting movies-app at 3000..")
	fmt.Println("initializing mongo connection..")
	helpers.ConnectDB()
	router := mux.NewRouter()

	router.HandleFunc("/api/movie/search/{name}", api.PublicSearch).Methods("GET")
	router.Handle("/api/movie/list", middleware.Login(http.HandlerFunc(api.SearchByUserId))).Methods("GET")
	router.Handle("/api/movie", middleware.AdminLogin(http.HandlerFunc(api.CreateMovie))).Methods("POST")
	router.Handle("/api/movie/rating/{id}", middleware.Login(http.HandlerFunc(api.RateMovie))).Methods("PUT")
	router.Handle("/api/movie/comment/{id}", middleware.Login(http.HandlerFunc(api.CommentMovie))).Methods("PUT")

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		fmt.Println(err)
	}
}
