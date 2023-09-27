package main

import (
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
	"github.com/charlesozo/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	stringChirpId := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(stringChirpId)
	if err!=nil{
	  respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
	  return
	}

    token, err := auth.GetBearerToken(r.Header)
	if err != nil {
      respondWithError(w, http.StatusUnauthorized, "Unauthorized to Perform Action")
	  return 
	}
	subject, err := auth.ValidateJWT(token, cfg.jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt validate JWT")
		return 
	}

	userID, err := strconv.Atoi(subject)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldnt parse userID")
		return 
	}
	dbChirp, err := cfg.DB.GetChirpID(chirpID)
	if err!=nil{
		respondWithError(w, http.StatusNotFound, "Couldn't get Chirp")
		return
	}
	if dbChirp.AuthorID != userID{
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp")
		return
	}
	err = cfg.DB.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})

}