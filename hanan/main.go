package main

import (
	"encoding/json"
	"net/http"
)

type user struct {
	Username  string
	Firstname string
	Lastname  string
	Age       int
	Gender    string
	Email     string
}

const loggedIn = false

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/profile", profileHandler)

	http.ListenAndServe(":2004", mux)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	if r.Method != "GET" {
		http.Error(w, "تباً لك، GET فقط. بنق", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Query().Get("token") != "okok" {
		http.Error(w, "You are not logged in", http.StatusUnauthorized)
		return
	}

	// get from DB
	user := user{
		Username:  "Hanan",
		Firstname: "Budala",
		Lastname:  "Hanin",
		Email:     "bng",
	}

	myData, _ := json.Marshal(user)

	w.Write(myData)
}
