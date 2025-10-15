package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

const PORT = "8080"
const APP_ROOT = "/app"
const FILE_SERVER_ROOT = "."

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	handler := http.NewServeMux()
	apiConfig := apiConfig{}

	var routes = map[string]func(w http.ResponseWriter, r *http.Request){
		"GET /api/healthz":         handleOK,
		"POST /api/validate_chirp": validateChirp,
		"POST /admin/reset":        apiConfig.resetNumberOfHits,
		"GET /admin/metrics":       apiConfig.getNumberOfHits,
	}

	handler.Handle(APP_ROOT+"/", apiConfig.middlewareMetricsInc(
		http.StripPrefix(
			APP_ROOT, http.FileServer(http.Dir(FILE_SERVER_ROOT)))))
	for route, handleFunc := range routes {
		handler.HandleFunc(route, handleFunc)
	}

	server := http.Server{
		Addr:    ":" + PORT,
		Handler: middlewareLogReq(handler),
	}

	log.Printf("Serving files from %s on port: %s\n", FILE_SERVER_ROOT, PORT)
	log.Fatal(server.ListenAndServe())
}
