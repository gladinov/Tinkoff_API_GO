package service

import (
	"errors"
	"math"
	"time"

	"github.com/gothanks/myapp/api/moex"
	"github.com/gothanks/myapp/other_func"
)

const (
	layout     = "2006-01-02"
	hoursInDay = 24
	daysInYear = 365
)

type Report struct {
	BondsInRUB []BondReport
	BondsInCNY []BondReport
}

type BondReport struct {
	Name                      string
	Ticker                    string
	MaturityDate              string // Дата погашения
	OfferDate                 string
	Duration                  int64
	BuyDate                   string
	BuyPrice                  float64
	YieldToMaturityOnPurchase float64 // Доходность к погашению при покупке
	YieldToOfferOnPurchase    float64 // Доходность к оферте при покупке
	YieldToMaturity           float64 // Текущая доходность к погашению
	YieldToOffer              float64 // Текущая доходность к оферте
	CurrentPrice              float64
	Nominal                   float64
	Profit                    float64 // Результат инвестиции
	AnnualizedReturn          float64 // Годовая доходность
}

func (resultReports *Report) CreateBondReport(reportPostions ReportPositions) error {

	for i := range reportPostions.CurrentPositions {
		position := reportPostions.CurrentPositions[i]
		switch position.Currency {
		case "rub":
			bondReport, err := createBondReport(position)
			if err != nil {
				return errors.New("service: GetBondReport" + err.Error())
			}
			resultReports.BondsInRUB = append(resultReports.BondsInRUB, bondReport)
		case "cny":
			bondReport, err := createBondReport(position)
			if err != nil {
				return errors.New("service: GetBondReport" + err.Error())
			}
			resultReports.BondsInCNY = append(resultReports.BondsInCNY, bondReport)
		default:
			continue
		}
	}
	return nil
}

func createBondReport(position SharePosition) (BondReport, error) {

	var bondReport BondReport
	var moexBuyData moex.MoexUnmarshallStruct
	var moexLastPriceData moex.MoexUnmarshallStruct
	err := moexBuyData.GetSpecifications(position.Ticker, position.BuyDate)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}
	date := time.Now()
	err = moexLastPriceData.GetSpecifications(position.Ticker, date)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}

	bondReport = BondReport{
		Name:         position.Name,
		Ticker:       position.Ticker,
		BuyDate:      position.BuyDate.Format(layout),
		BuyPrice:     other_func.RoundFloat(position.BuyPrice, 2),
		CurrentPrice: other_func.RoundFloat(position.SellPrice, 2),
		Nominal:      position.Nominal,
	}
	// Протестить и переписать
	if len(moexLastPriceData.Yields.History.Data) != 0 {
		maturityDate := moexLastPriceData.Yields.History.Data[0].MaturityDate
		if maturityDate != nil {
			bondReport.MaturityDate = *maturityDate
		}

		offerDate := moexLastPriceData.Yields.History.Data[0].OfferDate
		if offerDate != nil {
			bondReport.OfferDate = *offerDate
		}

		duration := moexLastPriceData.Yields.History.Data[0].Duration
		if duration != nil {
			bondReport.Duration = int64(*duration)
		}

		yieldToMaturity := moexLastPriceData.Yields.History.Data[0].YieldToMaturity
		if yieldToMaturity != nil {
			bondReport.YieldToMaturity = other_func.RoundFloat(*yieldToMaturity, 2)
		}

		yieldToOffer := moexLastPriceData.Yields.History.Data[0].YieldToOffer
		if yieldToOffer != nil {
			bondReport.YieldToOffer = other_func.RoundFloat(*yieldToOffer, 2)
		}
	}
	// Протестить и переписать
	if len(moexLastPriceData.Yields.History.Data) != 0 {
		yieldToMaturityOnPurchase := moexBuyData.Yields.History.Data[0].YieldToMaturity
		if yieldToMaturityOnPurchase != nil {
			bondReport.YieldToMaturityOnPurchase = other_func.RoundFloat(*yieldToMaturityOnPurchase, 2)
		}

		yieldToOfferOnPurchase := moexBuyData.Yields.History.Data[0].YieldToOffer
		if yieldToOfferOnPurchase != nil {
			bondReport.YieldToOfferOnPurchase = other_func.RoundFloat(*yieldToOfferOnPurchase, 2)
		}

	}
	profitInPercentage, err := getProfit(position)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}
	bondReport.Profit = profitInPercentage

	annualizedReturn, err := getAnnualizedReturnInPercentage(position)
	if err != nil {
		return bondReport, errors.New("service: createBondReport" + err.Error())
	}
	bondReport.AnnualizedReturn = annualizedReturn

	return bondReport, nil
}

func getProfit(p SharePosition) (float64, error) {
	profitWithoutTax := getSecurityIncomeWithoutTax(p)
	totalTax := getTotalTaxFromPosition(p, profitWithoutTax)
	profit := getSecurityIncome(profitWithoutTax, totalTax)
	profitInPercentage, err := getProfitInPercentage(profit, p.BuyPrice, p.Quantity)
	if err != nil {
		return profitInPercentage, errors.New("service: GetProfit" + err.Error())
	}
	return profitInPercentage, nil
}

func getAnnualizedReturnInPercentage(p SharePosition) (float64, error) {
	var totalReturn float64
	profitWithoutTax := getSecurityIncomeWithoutTax(p)
	totalTax := getTotalTaxFromPosition(p, profitWithoutTax)
	profit := getSecurityIncome(profitWithoutTax, totalTax)
	buyDate := p.BuyDate
	// Костыль. Надо переписать когда-нибудь, так как для закрытх позиций данная функция работать не будет
	sellDate := time.Now()
	timeDurationInDays := sellDate.Sub(buyDate).Hours() / float64(hoursInDay)
	// Если покупка и продажа были совершены в один день, то берем минимум один день
	timeDurationInYears := math.Max(1, timeDurationInDays) / float64(daysInYear)
	if p.BuyPrice != 0 || p.Quantity != 0 {
		totalReturn = profit / (p.BuyPrice * p.Quantity)
	} else {
		return 0, errors.New("service: getAnnualizedReturn : divide by zero")
	}
	annualizedReturn := math.Pow((1+totalReturn), (1/timeDurationInYears)) - 1
	annualizedReturnInPercentage := annualizedReturn * 100
	annualizedReturnInPercentageRound := other_func.RoundFloat(annualizedReturnInPercentage, 2)

	return annualizedReturnInPercentageRound, nil

}
