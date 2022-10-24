package main

import (
	"net/http"
	"os"
)

const (
	SERVER_URL = "http://localhost:8080/cotacao"
)

func main() {
	_, err := http.Get(SERVER_URL)
	if err != nil {
		panic(err)
	}
}

func CreateFile() {
	file, err := os.Create()
	if err != nil {
		panic(err)
	}
	defer file.Close()
}
