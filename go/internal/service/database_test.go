package service_test

import (
	"database/sql"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/summervik/swing-ranger/internal/model"
	"github.com/summervik/swing-ranger/internal/service"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type secrets struct {
	ConnectionStrings map[string]string `json:"ConnectionStrings"`
}

func TestUpsertAndGetEodPrice(t *testing.T) {
	data, err := os.ReadFile("../../testdata/secrets.json")
	if err != nil {
		t.Fatal(err)
	}
	var s secrets
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}

	connStr, ok := s.ConnectionStrings["Command"]
	if !ok {
		t.Fatal("Command connection string not found in secrets.json")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	p := model.EodPrice{
		Symbol:  "TESTCOMP",
		DateEod: time.Date(2025, 5, 17, 0, 0, 0, 0, time.UTC),
		Open:    decimal.NewFromFloat(150.25),
		High:    decimal.NewFromFloat(152.30),
		Low:     decimal.NewFromFloat(149.80),
		Close:   decimal.NewFromFloat(151.75),
		Volume:  45230000.0,
	}

	by := "test-user"

	err = service.UpsertEodPrices(db, []model.EodPrice{p}, by)
	if err != nil {
		t.Fatal(err)
	}

	prices, err := service.GetEodPrices(db, p.Symbol)
	if err != nil {
		t.Fatal(err)
	}

	if len(prices) == 0 {
		t.Fatal("expected at least one record")
	}

	got := prices[len(prices)-1]

	if got.Symbol != p.Symbol ||
		!got.DateEod.Equal(p.DateEod) ||
		!got.Open.Equal(p.Open) ||
		!got.High.Equal(p.High) ||
		!got.Low.Equal(p.Low) ||
		!got.Close.Equal(p.Close) ||
		got.Volume != p.Volume {
		t.Fatalf("data mismatch\ngot: %+v\nwant: %+v", got, p)
	}
}
