package main

import (
	"log"
	"net/http"
	"slices"
	"strings"
)

var PROFANE = []string{"kerfuffle", "sharbert", "fornax"}
const CENSORED = "****"

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body *string `json:"body"`
	}
	type SanitizedChirp struct {
		CleanedBody *string `json:"cleaned_body"`
	}

	chirp := Chirp{}
	if err := decodeJSON(r, &chirp); err != nil {
		log.Printf("Error decoding JSON: %s", err)
		writeErrorJSON(w, http.StatusInternalServerError, "Something went wrong")
	} else if chirp.Body == nil || *chirp.Body == "" {
		writeErrorJSON(w, http.StatusBadRequest, "Chirp body is required")
	} else if len(*chirp.Body) > 140 {
		writeErrorJSON(w, http.StatusBadRequest, "Chirp is too long")
	} else {
		writeJSON(w, http.StatusOK, SanitizedChirp{
			CleanedBody: sanitizeChirp(*chirp.Body),
		})
	}
}

func sanitizeChirp(chirp string) *string {
	sanitized := []string{}
	for word := range strings.SplitSeq(chirp, " ") {
		if slices.Contains(PROFANE, strings.ToLower(word)) {
			sanitized = append(sanitized, CENSORED)
		} else {
			sanitized = append(sanitized, word)
		}
	}
	
	result := strings.Join(sanitized, " ")
	return &result
}