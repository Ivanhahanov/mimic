package ui

import (
	"fmt"      // formatting and printing values to the console.
	"log"      // logging messages to the console.
	"net/http" // Used for build HTTP servers and clients.
	"strings"
	"time"
)

// Handler functions.
func (e *Endpoint) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, e.Resp)
}

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
	})
}

type IngressHTTP struct {
	Port      int32      `json:"port"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	URL  string `json:"url"`
	Resp string `json:"resp"`
}

func (s *IngressHTTP) Run() {
	mux := http.NewServeMux()
	// Registering our handler functions, and creating paths.
	fmt.Println(strings.Repeat("=", 20))
	for _, e := range s.Endpoints {
		mux.HandleFunc(e.URL, e.Handler)
		fmt.Printf(" %s => %s\n", e.URL, e.Resp)
	}
	fmt.Println(strings.Repeat("=", 20))

	log.Println("Started server on port", s.Port)
	logger := logRequestMiddleware(mux)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Handler:      logger,
	}

	// Spinning up the server.
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
