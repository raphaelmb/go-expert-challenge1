package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	SERVER_URL = "http://localhost:8080/cotacao"
)

func main() {
	c := GetCotacaoDolar(SERVER_URL)
	fmt.Println(string(c))
}

func GetCotacaoDolar(url string) []byte {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

// TODO: create file
func CreateFile() {
}
