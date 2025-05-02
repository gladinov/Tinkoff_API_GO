package tinkoffApi

import (
	"errors"

	"github.com/gothanks/myapp/other_func"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

type BondIdentIdentifiers struct {
	Ticker    string
	ClassCode string
	Name      string
	Nominal   float64
}

func GetBondsActionsFromTinkoff(client *investgo.Client, instrumentUid string) (BondIdentIdentifiers, error) {
	var res BondIdentIdentifiers
	instrumentService := client.NewInstrumentsServiceClient()
	bondUid, err := instrumentService.BondByUid(instrumentUid)
	if err != nil {
		return res, errors.New("GetTickerFromUid: instrumentService.BondByUid" + err.Error())
	}
	res.Ticker = bondUid.BondResponse.Instrument.GetTicker()
	res.ClassCode = bondUid.BondResponse.Instrument.GetClassCode()
	res.Name = bondUid.BondResponse.Instrument.GetName()
	res.Nominal = other_func.MoneyValue(bondUid.BondResponse.Instrument.GetNominal())

	return res, nil
}

func GetLastPriceFromTinkoffInPersentageToNominal(client *investgo.Client, instrumentUid string) (float64, error) {
	marketDataClient := client.NewMarketDataServiceClient()
	lastPriceAnswer, err := marketDataClient.GetLastPrices([]string{instrumentUid})
	if err != nil {
		return 0, errors.New("tinkoffApi:GetLastPriceFromTinkoff" + err.Error())
	}
	if len(lastPriceAnswer.LastPrices) == 0 {
		return 0, errors.New("tinkoffApi:GetLastPriceFromTinkoff: no price data")
	}

	lastPrice := lastPriceAnswer.LastPrices[0].Price.ToFloat()

	return lastPrice, nil
}
