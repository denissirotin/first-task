package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := DB.Create(&message).Error; err != nil {
		http.Error(w, "Could not save message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	var messages []Message
	if err := DB.Find(&messages).Error; err != nil {
		http.Error(w, "Could not retrieve messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func GetMessageByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var message Message
	if err := DB.First(&message, id).Error; err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var message Message
	if err := DB.First(&message, id).Error; err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := DB.Save(&message).Error; err != nil {
		http.Error(w, "Could not update message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := DB.Delete(&Message{}, id).Error; err != nil {
		http.Error(w, "Could not delete message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	InitDB()

	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	router.HandleFunc("/api/messages/{id}", GetMessageByID).Methods("GET")
	router.HandleFunc("/api/messages/{id}", UpdateMessage).Methods("PATCH")
	router.HandleFunc("/api/messages/{id}", DeleteMessage).Methods("DELETE")

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", router)
}
