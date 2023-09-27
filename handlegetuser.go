package main

import (
	"encoding/json"
	"net/http"
	"errors"
	"github.com/charlesozo/chirpy/internal/auth"
	"github.com/charlesozo/chirpy/internal/database"
)
type User struct{
	Email string `json:"email"`
	ID int `json:"id"`
    Password string `json:"-"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}
func (cfg *apiConfig) handlerChirpUsers(w http.ResponseWriter, r *http.Request){
    type parameters struct{
	    Email string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	} 
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Error decoding message")
	}
    hashedpassword, err := auth.HashPassword(params.Password)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}
	user,err := cfg.DB.CreateUser(params.Email, hashedpassword, false)
	if err != nil{
		if errors.Is(err, database.ErrAlreadyExists) {
			respondWithError(w, http.StatusConflict, "User already exists")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}
	respondWithJSON(w, http.StatusCreated, response{
		User: User{
         ID: user.ID,
		 Email: user.Email,
		 IsChirpyRed: user.IsChirpyRed,
		},
	})
}