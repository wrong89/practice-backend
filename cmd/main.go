package main

import (
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {}

func main() {
	http.HandleFunc("/test", handler)

	http.ListenAndServe(":9091", nil)
}
