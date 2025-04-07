package main

import (
	"errors"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

type Bond struct {
	Identifiers              Identifiers
	Name                     string         // GetBondsActionsFromPortfolio
	InstrumentType           string         // T_Api_Getportfolio
	Currency                 string         // T_Api_Getportfolio
	Quantity                 float64        // T_Api_Getportfolio
	AveragePositionPrice     float64        // T_Api_Getportfolio
	ExpectedYield            float64        // T_Api_Getportfolio
	CurrentNkd               float64        // T_Api_Getportfolio
	CurrentPrice             float64        // T_Api_Getportfolio
	AveragePositionPriceFifo float64        // T_Api_Getportfolio
	Blocked                  bool           // T_Api_Getportfolio
	ExpectedYieldFifo        float64        // T_Api_Getportfolio
	DailyYield               float64        // T_Api_Getportfolio
	Amortizations            *Amortizations // GetBondsActionsFromPortfolio
	Coupons                  *Coupons       // GetBondsActionsFromPortfolio
	Offers                   *Offers        // GetBondsActionsFromPortfolio
	Duration                 Duration       // GetBondsActionsFromPortfolio
}

type Identifiers struct {
	Ticker        string // GetBondsActionsFromPortfolio
	ClassCode     string // GetBondsActionsFromPortfolio
	Figi          string // T_Api_Getportfolio
	InstrumentUid string // T_Api_Getportfolio
	PositionUid   string // T_Api_Getportfolio
	AssetUid      string // GetBondsActionsFromPortfolio
}

func (b *Bond) GetBondsActionsFromPortfolio(client *investgo.Client) error {
	instrumentService := client.NewInstrumentsServiceClient()
	bondUid, err := instrumentService.BondByUid(b.Identifiers.InstrumentUid)
	if err != nil {
		return errors.New("GetTickerFromUid: instrumentService.BondByUid" + err.Error())
	}
	b.Identifiers.Ticker = bondUid.BondResponse.Instrument.GetTicker()
	b.Identifiers.ClassCode = bondUid.BondResponse.Instrument.GetClassCode()
	b.Name = bondUid.BondResponse.Instrument.GetName()
	return nil
}

// Получение данных
func (b *Bond) GetActionFromMoex() error {
	MoexUnmarshallData := MoexUnmarshallStruct{}
	err := MoexUnmarshallData.GetBondsFromMOEX(b.Identifiers.Ticker, 0, 20)
	if err != nil {
		return errors.New("GetBondsActionsFromPortfolio: GetBondsFromMOEX" + err.Error())
	}
	err = MoexUnmarshallData.GetDurationFromMoex(b.Identifiers.Ticker, b.Identifiers.ClassCode)
	if err != nil {
		return errors.New("GetBondsActionsFromPortfolio: GetBondsFromMOEX" + err.Error())
	}
	b.Amortizations = MoexUnmarshallData.Amortizations
	b.Offers = MoexUnmarshallData.Offers
	b.Coupons = MoexUnmarshallData.Coupons
	b.Duration = MoexUnmarshallData.Duration
	return nil
}
