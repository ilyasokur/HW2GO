package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const version = "0.0.1"

var healthy int32

type DecodeRequest struct {
	InputString string `json:"inputString"`
}

type DecodeResponse struct {
	OutputString string `json:"outputString"`
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"version":"%s"}`, version)
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	var req DecodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(req.InputString)
	if err != nil {
		http.Error(w, "Invalid base64 string", http.StatusBadRequest)
		return
	}

	resp := DecodeResponse{OutputString: string(decoded)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func hardOpHandler(w http.ResponseWriter, r *http.Request) {
	sleepDuration := time.Duration(10+rand.Intn(10)) * time.Second
	time.Sleep(sleepDuration)
	if rand.Intn(2) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/version", versionHandler)
	mux.HandleFunc("/decode", decodeHandler)
	mux.HandleFunc("/hard-op", hardOpHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Println("Server is starting on port 8080...")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}

	<-idleConnsClosed
	log.Println("Server gracefully stopped")
}
