package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)
func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
  stringChirpId := chi.URLParam(r, "chirpID")
  chirpID, err := strconv.Atoi(stringChirpId)
  if err!=nil{
	respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
  }
  dbChirp, err := cfg.DB.GetChirpID(chirpID)
  if err !=nil{
	respondWithError(w, http.StatusNotFound, err.Error())
  }
  respondWithJSON(w, http.StatusOK, Chirp{
	ID: dbChirp.ID,
	Body: dbChirp.Body,
  })
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			AuthorID: dbChirp.AuthorID,
			Body: dbChirp.Body,
			ID:   dbChirp.ID,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}