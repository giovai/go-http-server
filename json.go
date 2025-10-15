package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func writeErrorJSON(w http.ResponseWriter, httpCode int, errorMessage string) {
	type error struct {
		Error string `json:"error"`
	}
	writeJSON(w, httpCode, error{
		Error: errorMessage,
	})
}

func writeJSON(w http.ResponseWriter, httpCode int, data any) {
	responseJSON, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error marshalling JSON: %s", err)
		return
	}
	w.WriteHeader(httpCode)
	w.Header().Set("Content-Type", "application/json")
	log.Printf("%v", w.Header().Values("Content-Type"))
	w.Write(responseJSON)
}

func decodeJSON(r *http.Request, data any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&data)
}
