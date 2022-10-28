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

var DB *sql.DB

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
	db, _ := sql.Open("sqlite3", "db.sqlite")
	defer db.Close()
	DB = db
	prepareAndConnectDB()
	http.HandleFunc("/cotacao", CotacaoHandler)
	http.ListenAndServe(PORT, nil)
}

func prepareAndConnectDB() {
	os.Remove("db.sqlite")
	file, err := os.Create("db.sqlite")
	if err != nil {
		panic(err)
	}
	file.Close()

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS cotacao (
		code TEXT, code_in TEXT, name TEXT, high TEXT, low TEXT, var_bid TEXT, pct_change TEXT, bid TEXT, ask TEXT, timestamp TEXT, create_date TEXT
		);
	`)
	if err != nil {
		panic(err)
	}
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	// Change ctx time
	case <-time.After(time.Second * 10):
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

	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
	}
}

func GetCotacao(url string) (*Cotacao, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	err = SaveCotacao(DB, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// TODO: change ctx time
func SaveCotacao(db *sql.DB, cotacao *Cotacao) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, "insert into cotacao(code, code_in, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, cotacao.USDBRL.Code, cotacao.USDBRL.Codein, cotacao.USDBRL.Name, cotacao.USDBRL.High, cotacao.USDBRL.Low, cotacao.USDBRL.VarBid, cotacao.USDBRL.PctChange, cotacao.USDBRL.Bid, cotacao.USDBRL.Ask, cotacao.USDBRL.Timestamp, cotacao.USDBRL.CreateDate)
	if err != nil {
		return err
	}

	return nil
}
