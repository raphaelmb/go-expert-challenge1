package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// SERVER_URL = "http://localhost:8080/cotacao"
	SERVER_URL = "http://server:8080/cotacao"
)

func main() {
	c := GetCotacaoDolar(SERVER_URL)
	CreateAndWriteFile(string(c))
}

func GetCotacaoDolar(url string) []byte {
	ctx := context.Background()
	// Change context time
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func CreateAndWriteFile(cotacao string) {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write([]byte("Dólar: " + cotacao))
	if err != nil {
		panic(err)
	}
}
