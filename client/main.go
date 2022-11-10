package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	SERVER_URL = "http://server:8080/cotacao"
)

func main() {
	c, err := GetCotacaoDolar(SERVER_URL)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	CreateAndWriteFile(string(c))
	log.Println("Request completed successfully.")
}

func GetCotacaoDolar(url string) ([]byte, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
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
