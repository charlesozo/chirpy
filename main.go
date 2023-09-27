package main

import (
	"flag"
	"os"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/charlesozo/chirpy/internal/database"
	"github.com/joho/godotenv"
)
type apiConfig struct{
	fileserverHits int
	DB *database.DB
	jwtsecret string
	polkakey string
}

func main() {
	godotenv.Load(".env")
	const filepathRoot = "."
	const port = "8080"
	
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	if jwtSecret==""{
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	db, err :=  database.NewDB("database.json")
	if err !=nil{
		log.Fatal(err)
	}
    apiCfg  := apiConfig{
		fileserverHits: 0,
		DB: db,
		jwtsecret: jwtSecret,
		polkakey: polkaKey,

	}
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := apiCfg.DB.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}
	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

   	apiRouter := chi.NewRouter()
	webHookRouter := chi.NewRouter()

	webHookRouter.Post("/webhooks", apiCfg.handlerWebHooks)

	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handleReset)

	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsRetrieve)

	apiRouter.Delete("/chirps/{chirpID}", apiCfg.handlerChirpsDelete)
    
	apiRouter.Post("/users", apiCfg.handlerChirpUsers)
	apiRouter.Post("/refresh", apiCfg.handlerRefresh)
	apiRouter.Post("/revoke", apiCfg.handlerRevoke)
	apiRouter.Put("/users", apiCfg.handleUsersUpdate)
	apiRouter.Post("/login", apiCfg.handleLogin)

	apiRouter.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	apiRouter.Mount("/polka", webHookRouter)
	router.Mount("/api", apiRouter)
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handleMetrics)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}



