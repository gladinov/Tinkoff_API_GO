package moex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type MoexUnmarshallStruct struct {
	Amortizations *Amortizations `json:"amortizations"` // GetBondsActionsFromPortfolio
	Coupons       *Coupons       `json:"coupons"`       // GetBondsActionsFromPortfolio
	Offers        *Offers        `json:"offers"`        // GetBondsActionsFromPortfolio
	Duration      Duration
	Yields        Yields
}

// Получение Амортизации, купонов и офферов с MOEX
func (m *MoexUnmarshallStruct) GetBondsFromMOEX(ticker string, start, limit int) error {
	client := http.Client{Timeout: 3 * time.Second}
	uri := fmt.Sprintf("https://iss.moex.com/iss/securities/%s/bondization.json", ticker)

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

		var data MoexUnmarshallStruct
		err = json.Unmarshal(body, &data)
		if err != nil {
			return errors.New("GetBondsFromMoex: json.Unmarshall" + err.Error())
		}
		if len(data.Amortizations.Data) == 0 && len(data.Coupons.Data) == 0 && len(data.Offers.Data) == 0 {
			cont = false
		} else {
			if m.Amortizations == nil && m.Coupons == nil && m.Offers == nil {
				m.Amortizations = data.Amortizations
				m.Coupons = data.Coupons
				m.Offers = data.Offers
			} else {
				m.Amortizations.Data = append(m.Amortizations.Data, data.Amortizations.Data...)
				m.Coupons.Data = append(m.Coupons.Data, data.Coupons.Data...)
				m.Offers.Data = append(m.Offers.Data, data.Offers.Data...)
			}
			start += limit
		}

	}
	return nil
}

// Получение значения дюрации с MOEX
func (m *MoexUnmarshallStruct) GetDurationFromMoex(ticker string, class_code string) error {
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
	m.Duration = data.Marketdata_yields.Data

	return nil
}

func (m *MoexUnmarshallStruct) GetSpecifications(ticker string, date time.Time) error {
	formatDate := date.Format("2006-04-04")
	client := http.Client{Timeout: 3 * time.Second}
	uri := fmt.Sprintf("https://iss.moex.com/iss/history/engines/stock/markets/bonds/sessions/3/securities/%s.json", ticker)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return errors.New("GetSpecifications: request" + err.Error())
	}
	params := url.Values{}
	params.Add("limit", "1")
	params.Add("iss.meta", "off")
	params.Add("history.columns", "TRADEDATE,MATDATE,OFFERDATE,BUYBACKDATE,YIELDCLOSE,YIELDTOOFFER,FACEVALUE,DURATION")
	params.Add("limit", "1")
	params.Add("from", formatDate)
	params.Add("to", formatDate)

	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return errors.New("GetSpecifications: client.Do" + err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("GetSpecifications: resp.Body.Close" + err.Error())
	}

	var data Yields
	err = json.Unmarshal(body, &data)
	if err != nil {
		return errors.New("GetSpecifications: json.Unmarshal" + err.Error())
	}
	m.Yields = data
	return nil
}
