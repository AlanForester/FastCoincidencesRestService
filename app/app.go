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
	"strconv"
	"strings"
	"time"
)

var listenAddr = ":8080"

func main() {
	for SQL() == nil {
		log.Printf("Initializing DB ")
		time.Sleep(5 * time.Second)
	}
	// Loading test data
	if os.Getenv("APP_ENV") != "production" {
		var countRecords int
		var needRecords = 10000000
		SQL().Limit(1).Find(&mdl.ConnLog{}).Count(&countRecords)
		if countRecords == 0 {
			SQL().Find(&mdl.ConnLog{}).Count(&countRecords)
			log.Printf("CountRecords %d ", countRecords)
			if needRecords < countRecords {
				needRecords = needRecords - countRecords
				scripts.LoadData(needRecords)
			}
		} else {
			log.Printf("DB in not empty - OK")
		}
	}

	//mdl.LoadDuplicates()

	router := setupRouter()
	RunServer(router)
}

func setupRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", handleDups)
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
	//var logs []*mdl.ConnLog
	log.Printf("%v", r.URL.Path)
	params := strings.Split(r.URL.Path, "/")
	if len(params) == 3 {
		user1, user2 := 0, 0
		if user1t, err := strconv.Atoi(params[1]); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			user1 = user1t
		}
		if user2t, err := strconv.Atoi(params[2]); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			user2 = user2t
		}
		resp := mdl.IntersectionSQL(int64(user1), int64(user2))

		httpResp := false
		if len(resp) > 1 {
			httpResp = true
		}
		if json.NewEncoder(w).Encode(map[string]bool{"dupes": httpResp}) != nil {
			w.WriteHeader(http.StatusBadGateway)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	//user1 := r.URL.Query().Get(":id1")
	//user2 := r.URL.Query().Get(":id2")
	//if json.NewEncoder(w).Encode(user1,user2) != nil {
	//	w.WriteHeader(http.StatusBadGateway)
	//}
	//if json.NewEncoder(w).Encode(map[string]bool{"dupes": false}) != nil {
	//	w.WriteHeader(http.StatusBadGateway)
	//}
}
