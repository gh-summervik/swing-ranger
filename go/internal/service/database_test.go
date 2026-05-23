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

func TestUpsertAndGetEodCandlestick(t *testing.T) {
	data, err := os.ReadFile("../../testdata/secrets.json")
	if err != nil {
		t.Fatal(err)
	}
	var s secrets
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}

	cmdConnStr, ok := s.ConnectionStrings["Command"]
	if !ok {
		t.Fatal("Command connection string not found in secrets.json")
	}

	qryConnStr, ok := s.ConnectionStrings["Query"]
	if !ok {
		t.Fatal("Query connection string not found in secrets.json")
	}

	cmdDb, err := sql.Open("postgres", cmdConnStr)
	if err != nil {
		t.Fatal(err)
	}
	qryDb, err := sql.Open("postgres", qryConnStr)
	if err != nil {
		t.Fatal(err)
	}
	defer cmdDb.Close()
	defer qryDb.Close()

	dbService := &service.DbService{
		Command: cmdDb,
		Query:   qryDb,
	}

	c := model.NewEodCandlestick(
		"TESTCOMP",
		time.Date(2025, 5, 17, 0, 0, 0, 0, time.UTC),
		decimal.NewFromFloat(150.25),
		decimal.NewFromFloat(152.30),
		decimal.NewFromFloat(149.80),
		decimal.NewFromFloat(151.75),
		45230000.0,
	)

	by := "test-user"

	err = dbService.UpsertEodPrices([]model.EodCandlestick{c}, by)
	if err != nil {
		t.Fatal(err)
	}

	candles, err := dbService.GetEodCandlesticks(c.Symbol)
	if err != nil {
		t.Fatal(err)
	}

	if len(candles) == 0 {
		t.Fatal("expected at least one record")
	}

	got := candles[len(candles)-1]

	if got.Symbol != c.Symbol ||
		!got.DateEod.Equal(c.DateEod) ||
		!got.Open.Equal(c.Open) ||
		!got.High.Equal(c.High) ||
		!got.Low.Equal(c.Low) ||
		!got.Close.Equal(c.Close) ||
		got.Volume != c.Volume {
		t.Fatalf("data mismatch\ngot: %+v\nwant: %+v", got, c)
	}
}
