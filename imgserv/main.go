package main

import (
	"net/http"
	"time"

	"github.com/soffokl/go-progprobs/imgserv/handler"
)

func main() {
	img := &handler.Image{}
	http.Handle("/generate/", http.TimeoutHandler(img, time.Second*30, "Request timed out"))
	http.HandleFunc("/stats", handler.Stats)

	http.ListenAndServe(":8080", nil)
}
