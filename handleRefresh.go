package main
import (
	"net/http"
	"github.com/charlesozo/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
   type response struct{
	  Token string `json:"token"`
   }
  refreshToken, err := auth.GetBearerToken(r.Header)
  if err!=nil{
	respondWithError(w, http.StatusBadRequest, "Couldnt find JWT")
	return
  }
  isRevoked, err := cfg.DB.IsTokenRevoked(refreshToken)
  if err!=nil{
	    respondWithError(w, http.StatusInternalServerError, "Couldn't check session")
		return
  }
  if isRevoked{
	respondWithError(w, http.StatusUnauthorized, "Token is Revoke")
	return
  }
  accessToken, err := auth.RefreshToken(refreshToken, cfg.jwtsecret)
  if err != nil{
	respondWithError(w, http.StatusInternalServerError, "Couldn't refresh token")
  }
  respondWithJSON(w, http.StatusOK, response{
	Token: accessToken,
   })
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request){
     refreshToken, err := auth.GetBearerToken(r.Header)
	 if err != nil{
		respondWithError(w, http.StatusBadRequest, "Unable to find JWT")
		return
	}
	err = cfg.DB.Revoke(refreshToken)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "Unable to revoke")
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}