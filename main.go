package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)
type apiConfig struct{
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"
    apiCfg  := apiConfig{
		fileserverHits: 0,
	}
	router := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)
   	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handleReset)
	
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



