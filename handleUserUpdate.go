package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/charlesozo/chirpy/internal/auth"
)
func (cfg *apiConfig) handleUsersUpdate(w http.ResponseWriter, r *http.Request){
	type Parameters struct {
        Email string `json:"email"`
		Password string `json:"password"`
	}
	type response struct{
		User
	}
    token, err := auth.GetBearerToken(r.Header)
    if err!=nil{
		respondWithError(w, http.StatusUnauthorized, "Couldnt Find JWT")
		return
	}
	subject, err :=  auth.ValidateJWT(token, cfg.jwtsecret)
	if err!=nil{
		respondWithError(w, http.StatusUnauthorized, "Couldn't Validate JWT")
	}
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err = decoder.Decode(&params)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode parameters")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "Couldnt hash password")
		return 
	}
   userIDInt, err := strconv.Atoi(subject)

   if err != nil {
	respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
	return
    }
	updatedUser, err := cfg.DB.UpdateUser(userIDInt, params.Email, hashedPassword)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User{
			Email: updatedUser.Email,
			ID: updatedUser.ID,
		},
	})
}