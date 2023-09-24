package main
import(
	"log"
	"net/http"
	"encoding/json"
)

func respondWithError(w http.ResponseWriter, code int, msg string){
	if code >499{
	  log.Printf("Responding with 5XX error: %s", msg)
	}
	w.WriteHeader(code)
  }
  func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	  w.Header().Set("content-type", "application/json")
	  data, err := json.Marshal(payload)
	  if err != nil{
		  log.Printf("Error marshalling JSON: %s", err)
		  w.WriteHeader(500)
		  return
	  }
	  w.WriteHeader(code)
	  w.Write(data)
  }