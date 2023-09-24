package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

)
type Chirp struct{
	Body string `json:"body"`
	ID   int    `json:"id"`
}
func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Body string `json:"body"`
	}
     decoder := json.NewDecoder(r.Body)
	 params := parameters{}
	 err := decoder.Decode(&params)
	 if err != nil {
		respondWithError(w, 400, "Error decoding message")
		return
	}
	cleaned, err := validateChirp(params.Body)
    if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	chirp, err := cfg.DB.CreateChirp(cleaned)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		Body: chirp.Body,
		ID:   chirp.ID,
	})

}
func validateChirp(body string) (string, error){
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
	   return "", errors.New("chirp is too long")
   }
   words := strings.Fields(body)
	for i, word := range words {
		if strings.ToLower(word)=="kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word)== "fornax"{
			words[i] = "****"
		}
	}
	cleaned_words := strings.Join(words, " ")
	return cleaned_words,  nil
}