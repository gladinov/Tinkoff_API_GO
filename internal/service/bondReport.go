package service

import (
	"time"
)

type BondReports struct {
	Data []BondReport
}

type BondReport struct {
	Name                      string
	MaturityDate              time.Time // Дата погашения
	OfferDate                 time.Time
	Duration                  int64
	Ticker                    string
	BuyDate                   time.Time
	BuyPrice                  float64
	YieldToMaturityOnPurchase float64 // Доходность к погашению при покупке
	YieldToOfferOnPurchase    float64 // Доходность к оферте при покупке
	YieldToMaturity           float64 // Текущая доходность к погашению
	YieldToOffer              float64 // Текущая доходность к оферте
	CurrentPrice              float64
	Profit                    float64 // Результат инвестиции
	AnnualizedReturn          float64 // Годовая доходность
}
