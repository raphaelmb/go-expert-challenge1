package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	API_URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	PORT    = ":8080"
)

type Cotacao struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", CotacaoHandler)
	http.ListenAndServe(PORT, nil)
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	c, err := GetCotacao(API_URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c.USDBRL.Bid)
}

func GetCotacao(url string) (*Cotacao, error) {
	ctx := context.Background()
	// Change context time
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

	var c Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// TODO: save to database
func SaveCotacao() {}
