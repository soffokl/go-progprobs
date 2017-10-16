package main

import (
	"net/http"

	"github.com/soffokl/go-progprobs/imgserv/handler"
)

func main() {
	http.HandleFunc("/generate/", handler.Image)
	http.HandleFunc("/stats", handler.Stats)

	http.ListenAndServe(":8080", nil)
}
