package main

import (
	"flag"
	"log"
	"net/http"

	"lab9_var7/task5_go_compute_service/server"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9001", "listen address")
	flag.Parse()

	srv := &http.Server{
		Addr:    *addr,
		Handler: server.NewMux(),
	}

	log.Printf("primeservice listening on %s", *addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
