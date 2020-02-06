package main

import (
	"app/mdl"
	"app/scripts"
	. "app/srv"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var listenAddr = ":8080"

func main() {

	// Loading test data
	if os.Getenv("APP_ENV") != "production" {
		var countRecords int
		var needRecords = 10000000
		SQL().Find(&mdl.ConnLog{}).Count(&countRecords)
		log.Printf("countRecords %d ", countRecords)
		if countRecords < needRecords {
			needRecords = needRecords - countRecords
			scripts.LoadData(needRecords)
		}
	}

	router := setupRouter()
	RunServer(router)
}

func setupRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", notFound)
	router.HandleFunc("/1", handleDups)
	return router
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			exTimeStart := time.Now().UnixNano()
			defer func() {
				exTime := float64(time.Now().UnixNano()-exTimeStart) / 1000000 // To milliseconds
				logger.Printf("[%s] %s - %s %s | %.2fms", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), exTime)
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func RunServer(router *http.ServeMux) {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      logging(logger)(router),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Database closing...")
		_ = SQL().Close()
		logger.Println("Server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()
	logger.Println("Server is ready to handle requests at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}
	<-done
	logger.Println("Server stopped")

}

func notFound(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}

func handleDups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	var logs []*mdl.ConnLog
	SQL().Find(&logs)

	if json.NewEncoder(w).Encode(logs) != nil {
		w.WriteHeader(http.StatusBadGateway)
	}
	//if json.NewEncoder(w).Encode(map[string]bool{"dupes": false}) != nil {
	//	w.WriteHeader(http.StatusBadGateway)
	//}
}
