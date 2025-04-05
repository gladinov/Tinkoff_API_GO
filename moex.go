package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

type Bond struct {
	Ticker        string
	Amortizations *Amortizations `json:"amortizations"`
	Coupons       *Coupons       `json:"coupons"`
	Offers        *Offers        `json:"offers"`
	Duration      Duration
}

type Amortizations struct {
	Columns []string       `json:"columns"`
	Data    []Amortization `json:"data"`
}

type Amortization struct {
	Isin             string  `json:"isin"`
	Name             string  `json:"name"`
	Issuevalue       float64 `json:"issuevalue"`
	Amortdate        string  `json:"amortdate"`
	Facevalue        float64 `json:"facevalue"`
	Initialfactvalue float64 `json:"initialfactvalue"`
	Faceuint         string  `json:"faceunit"`
	Valueprc         float64 `json:"valueprc"`
	Value            float64 `json:"value"`
	Value_rub        float64 `json:"value_rub"`
	Data_source      string  `json:"data_source"`
	Secid            string  `json:"secid"`
	Primary_boardid  string  `json:"primary_boardid"`
}

func (a *Amortization) UnmarshalJSON(data []byte) error {
	dataSlice := make([]any, 13)
	err := json.Unmarshal(data, &dataSlice)
	if err != nil {
		return errors.New("CustomFloat64: UnmarshalJSON: " + err.Error())
	}
	a.Isin = dataSlice[0].(string)
	a.Name = dataSlice[1].(string)
	a.Issuevalue = dataSlice[2].(float64)
	a.Amortdate = dataSlice[3].(string)
	a.Facevalue = dataSlice[4].(float64)
	a.Initialfactvalue = dataSlice[5].(float64)
	a.Faceuint = dataSlice[6].(string)
	a.Valueprc = dataSlice[7].(float64)
	a.Value = dataSlice[8].(float64)
	a.Value_rub = dataSlice[9].(float64)
	a.Data_source = dataSlice[10].(string)
	a.Secid = dataSlice[11].(string)
	a.Primary_boardid = dataSlice[12].(string)

	return nil

}

type Coupons struct {
	Columns []string `json:"columns"`
	Data    []Coupon `json:"data"`
}

type Coupon struct {
	Isin             string   `json:"isin"`
	Name             string   `json:"name"`
	Issuevalue       float64  `json:"issuevalue"`
	Coupondate       string   `json:"coupondate"`
	Recorddate       *string  `json:"recorddate"`
	Startdate        string   `json:"startdate"`
	Initialfacevalue float64  `json:"initialfacevalue"`
	Facevalue        float64  `json:"facevalue"`
	Faceunit         string   `json:"faceunit"`
	Value            *float64 `json:"value,omitempty"`
	Valueprc         *float64 `json:"valueprc,omitempty"`
	Value_rub        *float64 `json:"value_rub,omitempty"`
	Secid            string   `json:"secid"`
	Primary_boardid  string   `json:"primary_boardid"`
}

func (c *Coupon) UnmarshalJSON(data []byte) error {
	dataSlice := make([]any, 14)
	err := json.Unmarshal(data, &dataSlice)
	if err != nil {
		return errors.New("CustomFloat64: UnmarshalJSON: " + err.Error())
	}
	c.Isin = dataSlice[0].(string)
	c.Name = dataSlice[1].(string)
	c.Issuevalue = dataSlice[2].(float64)
	c.Coupondate = dataSlice[3].(string)
	c.Recorddate = checkStringNull(dataSlice[4])
	c.Startdate = dataSlice[5].(string)
	c.Initialfacevalue = dataSlice[6].(float64)
	c.Facevalue = dataSlice[7].(float64)
	c.Faceunit = dataSlice[8].(string)
	c.Value = checkFloa64Null(dataSlice[9])
	c.Valueprc = checkFloa64Null(dataSlice[10])
	c.Value_rub = checkFloa64Null(dataSlice[11])
	c.Secid = dataSlice[12].(string)
	c.Primary_boardid = dataSlice[13].(string)

	return nil
}

func checkFloa64Null(a any) *float64 {
	if FloatVal, ok := a.(float64); ok {
		return &FloatVal
	} else {
		return nil
	}
}

type Offers struct {
	Columns []string `json:"columns"`
	Data    []Offer  `json:"data"`
}

type Offer struct {
	Isin            string   `json:"isin"`
	Name            string   `json:"name"`
	Issuevalue      float64  `json:"issuevalue"`
	Offerdate       string   `json:"offerdate"`
	Offerdatestart  string   `json:"offerdatestart"`
	Offerdateend    string   `json:"offerdateend"`
	Facevalue       float64  `json:"facevalue"`
	Faceunit        string   `json:"faceunit"`
	Price           *float64 `json:"price"`
	Value           float64  `json:"value"`
	Agent           *string  `json:"agent,omitempty"`
	Offertype       string   `json:"offertype"`
	Secid           string   `json:"secid"`
	Primary_boardid string   `json:"primary_boardid"`
}

func (o *Offer) UnmarshalJSON(data []byte) error {
	dataSlice := make([]any, 14) // Количество элементов соответствует полям структуры
	err := json.Unmarshal(data, &dataSlice)
	if err != nil {
		return errors.New("Offer: UnmarshalJSON: " + err.Error())
	}

	o.Isin = dataSlice[0].(string)
	o.Name = dataSlice[1].(string)
	o.Issuevalue = dataSlice[2].(float64)
	o.Offerdate = dataSlice[3].(string)
	o.Offerdatestart = dataSlice[4].(string)
	o.Offerdateend = dataSlice[5].(string)
	o.Facevalue = dataSlice[6].(float64)
	o.Faceunit = dataSlice[7].(string)
	o.Price = checkFloa64Null(dataSlice[8])
	o.Value = dataSlice[9].(float64)
	o.Agent = checkStringNull(dataSlice[10])
	o.Offertype = dataSlice[11].(string)
	o.Secid = dataSlice[12].(string)
	o.Primary_boardid = dataSlice[13].(string)

	return nil
}

func checkStringNull(a any) *string {
	if StringVal, ok := a.(string); ok {
		return &StringVal
	} else {
		return nil
	}
}

type Durations struct {
	Marketdata_yields Marketdata_yields
}

type Marketdata_yields struct {
	Data Duration
}

type Duration struct {
	Duration float64
	Date     time.Time
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	dataSlice := make([][]any, 0)
	err := json.Unmarshal(data, &dataSlice)
	if err != nil {
		return errors.New("CustomFloat64: UnmarshalJSON: " + err.Error())
	}
	d.Date = time.Now()
	d.Duration = dataSlice[0][0].(float64)
	return nil
}

// Получение Амортизации, купонов и офферов с MOEX
func (b *Bond) GetBondsFromMOEX(ticker string, start, limit int) error {
	b.Ticker = ticker

	client := http.Client{Timeout: 3 * time.Second}
	uri := fmt.Sprintf("https://iss.moex.com/iss/securities/%s/bondization.json", b.Ticker)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return errors.New("GetBondsFromMoex: request" + err.Error())
	}
	cont := true
	for cont {
		params := url.Values{}
		params.Add("start", strconv.Itoa(start))
		params.Add("limit", strconv.Itoa(limit))
		req.URL.RawQuery = params.Encode()

		resp, err := client.Do(req)
		if err != nil {
			return errors.New("GetBondsFromMoex: client.Do" + err.Error())
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.New("GetBondsFromMoex: io.ReadAll" + err.Error())
		}

		var data Bond
		err = json.Unmarshal(body, &data)
		if err != nil {
			return errors.New("GetBondsFromMoex: json.Unmarshall" + err.Error())
		}
		if len(data.Amortizations.Data) == 0 && len(data.Coupons.Data) == 0 && len(data.Offers.Data) == 0 {
			cont = false
		} else {
			if b.Amortizations == nil && b.Coupons == nil && b.Offers == nil {
				b.Amortizations = data.Amortizations
				b.Coupons = data.Coupons
				b.Offers = data.Offers
			} else {
				b.Amortizations.Data = append(b.Amortizations.Data, data.Amortizations.Data...)
				b.Coupons.Data = append(b.Coupons.Data, data.Coupons.Data...)
				b.Offers.Data = append(b.Offers.Data, data.Offers.Data...)
			}
			start += limit
		}

	}
	return nil
}

// Получение значения дюрации с MOEX
func (b *Bond) GetDurationFromMoex(ticker string, class_code string) error {
	client := http.Client{Timeout: 3 * time.Second}
	uri := fmt.Sprintf("https://iss.moex.com/iss/engines/stock/markets/bonds/boards/%s/securities/%s/securities.json", class_code, ticker)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return errors.New("GetDurationFromMoex: request" + err.Error())
	}
	params := url.Values{}
	params.Add("iss.only", "marketdata_yields")
	params.Add("iss.meta", "off")
	params.Add("marketdata_yields.columns", "DURATION")
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return errors.New("GetDurationFromMoex: client.Do" + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("GetDurationFromMoex: io.ReadAll" + err.Error())
	}
	var data Durations
	err = json.Unmarshal(body, &data)
	if err != nil {
		return errors.New("GetDurationFromMoex: json.Unmarshall" + err.Error())
	}
	b.Duration = data.Marketdata_yields.Data

	return nil
}

func GetBondsActionsFromPortfolio(uidTickerMap map[string][]string) (map[string]Bond, error) {
	bondRubList := make(map[string]Bond)
	for _, value := range uidTickerMap {
		ticker := value[0]
		class_code := value[1]
		b := Bond{}
		err := b.GetBondsFromMOEX(ticker, 0, 20)
		if err != nil {
			return nil, errors.New("GetBondsActionsFromPortfolio: GetBondsFromMOEX" + err.Error())
		}
		err = b.GetDurationFromMoex(ticker, class_code)
		if err != nil {
			return nil, errors.New("GetBondsActionsFromPortfolio: GetBondsFromMOEX" + err.Error())
		}
		bondRubList[ticker] = b
	}
	return bondRubList, nil
}

// Получаем Uid из базы данных
func GetUidsFromSQL(nameDB string) ([]string, error) {
	db, err := sql.Open("sqlite3", nameDB)
	if err != nil {
		return nil, errors.New("GetUidsFromSQL: sql.Open" + err.Error())
	}

	defer db.Close()
	// Выьираем из портфелей рублевые облигации
	query := "select instrumentUid from portfolio_2007907898 where instrumentType == 'bond' and currency == 'rub' union select instrumentUid from portfolio_2016119489 where instrumentType == 'bond' and currency == 'rub'"
	rows, err := db.Query(query)
	if err != nil {
		return nil, errors.New("GetUidsFromSQL: db.Query" + err.Error())
	}

	defer rows.Close()

	portfolioUids := make([]string, 0)
	for rows.Next() {
		var uid string
		err := rows.Scan(&uid)
		if err != nil {
			return nil, errors.New("GetUidsFromSQL: rows.Scan" + err.Error())
		}
		portfolioUids = append(portfolioUids, uid)
	}
	if err := rows.Err(); err != nil { // (5)
		return nil, errors.New("GetUidsFromSQL: rows.Err()" + err.Error())
	}
	return portfolioUids, nil
}

// получение ticker и class_code по uid
func GetTickerFromUid(instrumentService *investgo.InstrumentsServiceClient,
	uid string) ([]string, error) {
	retList := make([]string, 0)
	bondUid, err := instrumentService.BondByUid(uid)
	if err != nil {
		return nil, errors.New("GetTickerFromUid: instrumentService.BondByUid" + err.Error())
	}
	ticker := bondUid.BondResponse.Instrument.GetTicker()
	class_code := bondUid.BondResponse.Instrument.GetClassCode()
	retList = append(retList, ticker, class_code)
	return retList, nil
}

// создание МАПЫ по ключу uid и значению [ticker, class_code]
func MatchTickerAndUid(instrumentService *investgo.InstrumentsServiceClient,
	uids []string) (map[string][]string, error) {
	uidsTickerMap := make(map[string][]string)
	var err error
	for _, v := range uids {
		uidsTickerMap[v], err = GetTickerFromUid(instrumentService, v)
		if err != nil {
			return nil, errors.New("MatchTickerAndUid: GetTickerFromUid" + err.Error())
		}
	}
	return uidsTickerMap, nil
}
