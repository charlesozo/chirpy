package main

import (
	"encoding/json"
	"net/http"
     "time"
	"github.com/charlesozo/chirpy/internal/auth"
)
func (cfg *apiConfig)handleLogin(w http.ResponseWriter, r *http.Request){
  		type parameters struct{
			Email string `json:"email"`
			Password string `json:"password"`
			ExpiresInSeconds int    `json:"expires_in_seconds"`
		} 
		type response struct{
			User
			Token string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
	    err := decoder.Decode(&params) 
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to decode")
		}
		user, err := cfg.DB.GetUserByEmail(params.Email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to get user")
			return
		}
		err = auth.ComparePassword(params.Password, user.HashedPassword)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid password")
			return
		}
		
		
		accessToken, err := auth.MakeJWT(
			user.ID,
			cfg.jwtsecret,
			time.Hour,
			auth.TokenTypeAccess,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
			return
		}
	
		refreshToken, err := auth.MakeJWT(
			user.ID,
			cfg.jwtsecret,
			time.Hour*24*30*6,
			auth.TokenTypeRefresh,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh JWT")
			return
		}

		respondWithJSON(w, http.StatusOK, response{
			User: User{
				ID:    user.ID,
				Email: user.Email,
				IsChirpyRed: user.IsChirpyRed,
			},
			Token: accessToken,
			RefreshToken: refreshToken,
		})
}