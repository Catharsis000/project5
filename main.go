package main

import (
	"github.com/Catharsis000/project5.git/service"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	srv := service.New()

	mux.HandleFunc("/vote", srv.Vote)

	mux.HandleFunc("/stats", srv.Stats)
	_ = http.ListenAndServe(":9000", mux)

}
