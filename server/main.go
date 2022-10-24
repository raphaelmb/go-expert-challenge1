package main

import "net/http"

const (
	API_URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	PORT    = ":8080"
)

func main() {
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {})
	http.ListenAndServe(PORT, nil)
}
