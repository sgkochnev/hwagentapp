package main

import (
	"context"
	"hwAgentApp/logger"
	"hwAgentApp/moex"
	"hwAgentApp/store/pg"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("./config/.env"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cl := moex.NewMOEXClient(http.DefaultClient)

	resSequrities, err := cl.ListOfSequrities(ctx,
		&moex.ConfigSecurity{
			Q:     "Sberbank",
			Lang:  moex.LangEN,
			Limit: 5,
			Start: 0,
		})
	if err != nil {
		logger.Log(err)
	}

	tFrom, _ := time.Parse("2006-01-02", "2021-05-17")
	tTell := time.Now()

	// f, _ := os.OpenFile("result.csv", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	// defer f.Close()

	var insert = `INSERT INTO history(boardid,tradedate,shortname,secid,numtrades,value,open,low,high,legalcloseprice,waprice,close,volume,marketprice2,marketprice3,admittedquote,mp2valtrd,marketprice3tradesvalue,admittedvalue,waval,tradingsession)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)`

	mx := &sync.Mutex{}
	db, err := pg.Dial()
	if err != nil {
		logger.Log(err)
	}
	wg := sync.WaitGroup{}

	for i := range resSequrities {
		wg.Add(1)
		go func(i int) {
			for {
				resHistory, err := cl.History(
					ctx,
					moex.EngineStock,
					moex.MarketShares,
					resSequrities[i].MarketpriceBoardid,
					resSequrities[i].Secid,
					&moex.ConfigHistory{
						From: tFrom.Format("2006-01-02"),
						Till: tTell.Format("2006-01-02"),
						Lang: moex.LangEN,
					},
				)
				if err != nil {
					logger.Log(err)
				}
				l := len(resHistory)
				if l == 0 {
					break
				}
				for _, v := range resHistory {
					// fmt.Fprintln(f, v)
					mx.Lock()
					db.Exec(context.TODO(),
						insert,
						v.Boardid, v.Tradedate, v.Shortname, v.Secid,
						v.Numtrades, v.Value, v.Open, v.Low,
						v.High, v.Legalcloseprice, v.Waprice, v.Close,
						v.Volume, v.Marketprice2, v.Marketprice3, v.Admittedquote,
						v.Mp2valtrd, v.Marketprice3tradesvalue, v.Admittedquote, v.Waval,
						v.Tradingsession,
					)
					mx.Unlock()
				}
				tFrom, err = time.Parse("2006-01-02", resHistory[l-1].Tradedate)
				if err != nil {
					logger.Log(err)
				}
				tFrom = tFrom.AddDate(0, 0, 1)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	if err := db.Close(context.TODO()); err != nil {
		logger.Log(err)
	}
}
