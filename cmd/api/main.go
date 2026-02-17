package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const address = "localhost:2000"

func main() {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	http.ListenAndServe(":3000", router)
}
