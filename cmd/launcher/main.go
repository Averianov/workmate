package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	openapi "workmate/api/v1"
)

const (
	certFile = "cert.pem"
	keyFile  = "key.pem"
	port     = ":8080"
)

func main() {
	log.Printf("WorkMate Task Manager Server starting...")

	tasksAPIService := openapi.NewTasksAPIService()
	tasksAPIController := openapi.NewTasksAPIController(tasksAPIService)

	router := openapi.NewRouter(tasksAPIController)

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	if !fileExists(certFile) || !fileExists(keyFile) {
		log.Printf("Starting HTTP server on %s", port)
		log.Fatal(server.ListenAndServe())
	} else {
		server.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			NextProtos: []string{"h2", "http/1.1"}, // Поддержка HTTP/2
		}
		log.Printf("Starting HTTPS/HTTP2 server on %s", port)
		log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	}

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
