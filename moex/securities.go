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

type ConfigSecurity struct {
	// Поиск инструмента по части Кода, Названию, ISIN, Идентификатору Эмитента, Номеру гос.регистрации.
	// Например: https://iss.moex.com/iss/securities.xml?q=MOEX
	// Слова длиной менее трёх букв игнорируются. Если параметром передано два слова через пробел. То каждое должно быть длиной не менее трёх букв.
	Q string

	// Язык результата: ru или en
	Lang string

	// Номер строки (отсчет с нуля), с которой следует начать порцию возвращаемых данных (см. рук-во разработчика).
	// Получение ответа без данных означает, что указанное значение превышает число строк, возвращаемых запросом.
	Start int

	// Количество выводимых инструментов (5, 10, 20,100)
	Limit int

	// Группировать выводимый результат по полю. Доступны значения group и type.
	GroupBy string

	// Фильтровать по типам group или type.
	// Зависит от значения фильтра group_by.
	GroupByFilter string

	// Рынок.
	Engine string

	// Торгуется
	//  0 нет
	//  1 да
	IsTrading string

	// Возврашаемые колонки таблицы (secid,name,marketplace_boardid).
	SequritiesColumns string
}

func rawQuerySecurities(cfg *ConfigSecurity) string {
	q := url.Values{}
	// q.Add("iss.meta", "off")
	if cfg == nil {
		return q.Encode()
	}
	if cfg.Q != "" {
		q.Add("q", cfg.Q)
	}
	if cfg.Lang != "" {
		q.Add("lang", cfg.Lang)
	}
	if cfg.Engine != "" {
		q.Add("engine", cfg.Engine)
	}
	if cfg.Start != 0 {
		q.Add("start", strconv.Itoa(cfg.Start))
	}
	if cfg.Limit != 0 {
		q.Add("limit", strconv.Itoa(cfg.Limit))
	}
	if cfg.IsTrading != "" {
		q.Add("is_trading", cfg.IsTrading)
	}
	if cfg.SequritiesColumns != "" {
		q.Add("securities.columns", cfg.SequritiesColumns)
	}
	if cfg.GroupBy != "" {
		q.Add("group_by", cfg.GroupBy)
	}
	if cfg.GroupByFilter != "" {
		q.Add("group_by_filter", cfg.GroupByFilter)
	}

	return q.Encode()
}

type ModelSequrities struct {
	Id                 int    `csv:"id,omitempty"`
	Secid              string `csv:"secid,omitempty"`
	Shortname          string `csv:"shortname,omitempty"`
	Regnumber          string `csv:"regnumber,omitempty"`
	Name               string `csv:"name,omitempty"`
	Isin               string `csv:"isin,omitempty"`
	IsTraded           int    `csv:"is_traded,omitempty"`
	EmitentId          int    `csv:"emitent_id,omitempty"`
	EmitentTitle       string `csv:"emitent_title,omitempty"`
	EmitentInn         string `csv:"emitent_inn,omitempty"`
	EmitentOkpo        string `csv:"emitent_okpo,omitempty"`
	Gosreg             string `csv:"gosreg,omitempty"`
	Type               string `csv:"type,omitempty"`
	Group              string `csv:"group,omitempty"`
	PrimaryBoardid     string `csv:"primary_boardid,omitempty"`
	MarketpriceBoardid string `csv:"marketprice_boardid,omitempty"`
}

func parseListOfSequritiesInCSV(b []byte) ([]ModelSequrities, error) {
	res := []ModelSequrities{}
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

// Список бумаг торгуемых на московской бирже.
func (m *moex) ListOfSequrities(ctx context.Context, cfg *ConfigSecurity) ([]ModelSequrities, error) {

	url := compileURL(_baseURL, _securities) + _csv

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("can not create request: %v", err)
	}

	req.URL.RawQuery = rawQuerySecurities(cfg)
	logger.Log(req.URL)

	result, err := m.get(ctx, req)
	if err != nil {
		return nil, err
	}

	return parseListOfSequritiesInCSV(result)
}
