package moex

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"hwAgentApp/logger"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gocarina/gocsv"
)

const (
	SortASC  = "asc"
	SortDESC = "desc"

	TradingsessionMorning = "0"
	TradingsessionMain    = "1"
	TradingsessionEvening = "2"
	TradingsessionOutcome = "3"
)

type ConfigHistory struct {
	//Направление сортировки.
	// - "asc" - По возрастанию значения
	// - "desc" - По убыванию
	SortOrder string

	//Дата, начиная с которой необходимо начать выводить данные. Формат: ГГГГ-ММ-ДД.
	From string //Дата, начиная с которой необходимо начать выводить данные. Формат: ГГГГ-ММ-ДД.

	// Дата, до которой выводить данные. Формат: ГГГГ-ММ-ДД
	Till string // Дата, до которой выводить данные. Формат: ГГГГ-ММ-ДД

	//Минимальное количество сделок с бумагой.
	Numtrades int

	// 	Номер строки (отсчет с нуля), с которой следует начать порцию возвращаемых данных (см. рук-во разработчика).
	// Получение ответа без данных означает, что указанное значение превышает число строк, возвращаемых запросом.
	Start int

	//Язык результата: ru или en
	Lang string

	//Количество выводимых бумаг доступны значения (1, 5, 10, 20, 50, 100)
	Limit int

	//Поле, по которому сортируется ответ.
	SortColumn string

	// 	Показать данные только за необходимую сессию (только для фондового рынка)
	//   0 - Утренняя
	//   1 - Основная
	//   2 - Вечерняя
	//   3 - Итого
	Tradingsession string
}

func rawQueryHistory(cfg *ConfigHistory) string {
	q := url.Values{}
	if cfg == nil {
		return ""
	}
	if cfg.SortOrder != "" {
		q.Add("sort_order", cfg.SortOrder)
	}
	if cfg.From != "" {
		q.Add("from", cfg.From)
	}
	if cfg.Till != "" {
		q.Add("till", cfg.Till)
	}
	if cfg.Numtrades != 0 {
		q.Add("numtrades", strconv.Itoa(cfg.Numtrades))
	}
	if cfg.Start != 0 {
		q.Add("start", strconv.Itoa(cfg.Start))
	}
	if cfg.Lang != "" {
		q.Add("lang", cfg.Lang)
	}
	if cfg.Limit != 0 {
		q.Add("limit", strconv.Itoa(cfg.Limit))
	}
	if cfg.SortColumn != "" {
		q.Add("sort_column", cfg.SortColumn)
	}
	if cfg.Tradingsession != "" {
		q.Add("tradingsession", cfg.Tradingsession)
	}
	return q.Encode()
}

type ModelHistory struct {
	Boardid                 string  `csv:"BOARDID,omitempty"`
	Tradedate               string  `csv:"TRADEDATE,omitempty"`
	Shortname               string  `csv:"SHORTNAME,omitempty"`
	Secid                   string  `csv:"SECID,omitempty"`
	Numtrades               int     `csv:"NUMTRADES,omitempty"`
	Value                   float64 `csv:"VALUE,omitempty"`
	Open                    float64 `csv:"OPEN,omitempty"`
	Low                     float64 `csv:"LOW,omitempty"`
	High                    float64 `csv:"HIGH,omitempty"`
	Legalcloseprice         float64 `csv:"LEGALCLOSEPRICE,omitempty"`
	Waprice                 float64 `csv:"WAPRICE,omitempty"`
	Close                   float64 `csv:"CLOSE,omitempty"`
	Volume                  int     `csv:"VOLUME,omitempty"`
	Marketprice2            float64 `csv:"MARKETPRICE2,omitempty"`
	Marketprice3            float64 `csv:"MARKETPRICE3,omitempty"`
	Admittedquote           float64 `csv:"ADMITTEDQUOTE,omitempty"`
	Mp2valtrd               float64 `csv:"MP2VALTRD,omitempty"`
	Marketprice3tradesvalue float64 `csv:"MARKETPRICE3TRADESVALUE,omitempty"`
	Admittedvalue           float64 `csv:"ADMITTEDVALUE,omitempty"`
	Waval                   int     `csv:"WAVAL,omitempty"`
	Tradingsession          int     `csv:"TRADINGSESSION,omitempty"`
}

// Получить историю торгов для указанной бумаги на указанном режиме торгов за указанный интервал дат.
func (m *moex) History(ctx context.Context, engins, market, board, security string, cfg *ConfigHistory) ([]ModelHistory, error) {
	url := compileURL(
		_baseURL, _history,
		_engines, engins,
		_markets, market,
		_boards, board,
		_securities, security,
	) + _csv

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("can not create request: %v", err)
	}

	req.URL.RawQuery = rawQueryHistory(cfg)
	logger.Log(req.URL)

	result, err := m.get(ctx, req)
	if err != nil {
		return nil, err
	}

	return parseHistoryInCSV(result)
}

func parseHistoryInCSV(b []byte) ([]ModelHistory, error) {
	res := []ModelHistory{}
	reader := bufio.NewReader(bytes.NewReader(b))
	if _, _, err := reader.ReadLine(); err != nil {
		return nil, fmt.Errorf("cannot read data: %v", err)
	}
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = ';'
		return r
	})
	if err := gocsv.Unmarshal(reader, &res); err != nil {
		return nil, err
	}
	return res, nil
}
