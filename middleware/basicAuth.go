package middleware

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"movies-app/helpers"
	"movies-app/structs"
	"net/http"
)

func Login(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !checkUsernameAndPassword(user, pass, false) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
			w.WriteHeader(401)
			return
		}
		handler(w, r)
	}
}

func AdminLogin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !checkUsernameAndPassword(user, pass, true) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
			w.WriteHeader(401)
			return
		}
		handler(w, r)
	}
}

func checkUsernameAndPassword(username, password string, isAdminCheck bool) bool {
	if username == "" || password == "" {
		return false
	}
	client := helpers.GlobalMongoClient
	query := bson.M{
		"_id": username}
	c := client.Collection("user")
	ctx := context.TODO()
	var user structs.User
	c.FindOne(ctx, query).Decode(&user)
	if isAdminCheck {
		return username == user.ID && password == user.Password && user.IsAdmin
	}
	return username == user.ID && password == user.Password
}
