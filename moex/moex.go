package moex

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

const (
	_baseURL = "http://iss.moex.com/iss"
	_engines = "engines"
	_history = "history"
	_markets = "markets"
	_boards  = "boards"
	// _sessions   = "session"
	_securities = "securities"
	// _json       = ".json"
	// _xml        = ".xml"
	_csv = ".csv"
)

const (
	LangRU = "ru"
	LangEN = "en"

	GroupByType  = "type"
	GroupByGroup = "group"

	EngineStock         = "stock"        //Фондовый рынок и рынок депозитов
	EngineState         = "state"        //Рынок ГЦБ (размещение)
	EngineCurrency      = "currency"     //Валютный рынок
	EngineFutures       = "futures"      //Срочный рынок
	EngineCommodity     = "commodity"    //Товарный рынок
	EngineInterventions = "intervention" //Товарные интервенции
	EngineOffboard      = "offboard"     //ОТС-система
	EngineAgro          = "agro"         //Агро
	EngineOtc           = "otc"          //OTC Система

	MarketIndex         = "index"         //Индексы фондового рынка
	MarketShares        = "shares"        //Рынок акций
	MarketBonds         = "bonds"         //Рынок облигаций
	MarketNdm           = "ndm"           //Режим переговорных сделок
	MarketOtc           = "otc"           //ОТС
	MarketCcp           = "ccp"           //РЕПО с ЦК
	MarketDeposit       = "deposit"       //Депозиты с ЦК
	MarketRepo          = "repo"          //Рынок сделок РЕПО
	MarketQnv           = "qnv"           //Квал. инвесторы
	MarketMamc          = "mamc"          //Мультивалютный рынок смешанных активов
	MarketForeignshares = "foreignshares" //Иностранные ц.б.
	MarketForeignndm    = "foreignndm"    //Иностранные ц.б. РПС
	MarketMoexboard     = "moexboard"     //MOEX Board
	MarketGcc           = "gcc"           // РЕПО с ЦК с КСУ
	MarketCredit        = "credit"        //Рынок кредитов
	MarketStandard      = "standard"      //Standard
	MarketClassica      = "classica"      //Classica
)

type moex struct {
	client *http.Client
}

var once sync.Once
var moexClient *moex

func NewMOEXClient(client *http.Client) *moex {
	once.Do(
		func() {
			moexClient = &moex{client: client}
		},
	)
	return moexClient
}

func (c *moex) get(ctx context.Context, req *http.Request) ([]byte, error) {

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can not make request: %v", err)
	}
	defer resp.Body.Close()

	buf := &bytes.Buffer{}

	if _, err := io.Copy(buf, resp.Body); err != nil {
		return nil, fmt.Errorf("can not read response body: %v", err)
	}
	return buf.Bytes(), nil
}

func compileURL(baseUrl string, args ...string) string {
	return baseUrl + "/" + strings.Join(args, "/")
}
