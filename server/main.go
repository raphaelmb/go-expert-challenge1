package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
	prepareDb()
	db, _ := sql.Open("sqlite3", "db.sqlite")
	defer db.Close()

	http.HandleFunc("/cotacao", CotacaoHandler)
	http.ListenAndServe(PORT, nil)
}

func prepareDb() {
	os.Remove("db.sqlite")
	file, err := os.Create("db.sqlite")
	if err != nil {
		panic(err)
	}
	file.Close()
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	c, err := GetCotacao(API_URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d, err := strconv.ParseFloat(c.USDBRL.Bid, 64)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(d)
}

func GetCotacao(url string) (*Cotacao, error) {
	ctx := context.Background()
	// Change context time
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
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

	var c Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// TODO: save to database
// func SaveCotacao(db *sql.DB, cotacao *Cotacao) error {

// }
