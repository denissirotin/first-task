package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var message = "World"

type requestBody struct {
	Message string `json:"message"`
}

func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody requestBody

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message = reqBody.Message
	fmt.Fprintln(w, "Message updated successfully!")
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %s", message)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/hello", HelloHandler).Methods("GET")

	router.HandleFunc("/api/update-message", UpdateMessageHandler).Methods("POST")

	http.ListenAndServe(":8080", router)
}
