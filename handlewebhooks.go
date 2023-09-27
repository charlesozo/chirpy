package main
import (
	"net/http"
	"encoding/json"
	"errors"
	"github.com/charlesozo/chirpy/internal/database"
	"github.com/charlesozo/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerWebHooks(w http.ResponseWriter, r *http.Request){
	type WebHookParameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	apikey, err := auth.GetAPIKey(r.Header)
	if err!=nil{
		respondWithError(w, http.StatusUnauthorized, "Unable to get APIkey")
		return
	}

	if cfg.polkakey != apikey{
		respondWithError(w, http.StatusUnauthorized, "Api key doesnt match")
		return
	}
	decode := json.NewDecoder(r.Body)
	params := WebHookParameters{}
	err = decode.Decode(&params)
	if err!=nil{
		respondWithError(w, 400, "Error decoding message")
		return
	}
	if params.Event != "user.upgraded"{
		respondWithJSON(w, http.StatusOK, "Couldnt upgrade user")
		return
	}
	_, err = cfg.DB.UpgradeUser(params.Data.UserID, true)
	if err!=nil{
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user")
			return
		}
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})

}