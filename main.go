package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Функция для создания новой записи
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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// Функция для получения всех записей
func GetMessages(w http.ResponseWriter, r *http.Request) {
	var messages []Message
	if err := DB.Find(&messages).Error; err != nil {
		http.Error(w, "Could not retrieve messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func main() {
	InitDB() // Инициализация БД

	// Автоматическая миграция структуры Message
	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", router)
}
