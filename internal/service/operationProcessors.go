package service

import (
	"errors"
	"math"
	"time"

	"github.com/gothanks/myapp/api/tinkoffApi"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

const (
	WithholdingOfPersonalIncomeTaxOnCoupons        = 2     // 2	Удержание НДФЛ по купонам.
	WithholdingOfPersonalIncomeTaxOnDividends      = 8     // 8    Удержание налога по дивидендам.
	PartialRedemptionOfBonds                       = 10    // 10	Частичное погашение облигаций.
	PurchaseOfSecurities                           = 15    // 15	Покупка ЦБ.
	PurchaseOfSecuritiesWithACard                  = 16    // 16	Покупка ЦБ с карты.
	TransferOfSecuritiesFromAnotherDepository      = 17    // 17	Перевод ценных бумаг из другого депозитария.
	WithhouldingACommissionForTheTransaction       = 19    // 19	Удержание комиссии за операцию.
	PaymentOfDividends                             = 21    // 21	Выплата дивидендов.
	SaleOfSecurities                               = 22    // 22	Продажа ЦБ.
	PaymentOfCoupons                               = 23    // 23 Выплата купонов.
	StampDuty                                      = 47    // 47	Гербовый сбор.
	TransferOfSecuritiesFromIISToABrokerageAccount = 57    // 57   Перевод ценных бумаг с ИИС на Брокерский счет
	EuroTransBuyCost                               = 240   //Стоимость Евротранса при переводе из другого депозитария
	threeYearInHours                               = 26304 // Три года в часах
	baseTaxRate                                    = 0.13  // Налог с продажи ЦБ
)

type ReportPositions struct {
	Quantity         float64
	CurrentPositions []SharePosition
	ClosePositions   []SharePosition
}

type SharePosition struct {
	Name                  string
	BuyDate               time.Time
	SellDate              time.Time
	BuyOperationID        string
	SellOperationID       string
	Quantity              float64
	Figi                  string
	InstrumentType        string
	InstrumentUid         string
	Ticker                string
	ClassCode             string
	Nominal               float64
	BuyPrice              float64
	SellPrice             float64 // Для открытых позиций.Текущая цена с биржи
	BuyPayment            float64
	SellPayment           float64
	Currency              string
	BuyAccruedInt         float64 // НКД при покупке
	SellAccruedInt        float64
	PartialEarlyRepayment float64 // Частичное досрочное гашение
	TotalCoupon           float64
	TotalDividend         float64
	TotalComission        float64
	PaidTax               float64 // Фактически уплаченный налог(Часть налога будет уплачена в конце года, либо при выводе средств)
	TotalTax              float64 // Налог рассчитанный
	PositionProfit        float64 // С учетом рассчитанных налогов(TotalTax)
	ProfitInPercentage    float64 // В процентах строковая переменная
}

func (s *SharePosition) GetSpecificationsFromTinkoff(client *investgo.Client) error {
	resSpecFromTinkoff, err := tinkoffApi.GetBondsActionsFromTinkoff(client, s.InstrumentUid)
	if err != nil {
		return errors.New("service:GetSpecificationsFromMoex" + err.Error())
	}
	s.Ticker = resSpecFromTinkoff.Ticker
	s.ClassCode = resSpecFromTinkoff.ClassCode
	s.Nominal = resSpecFromTinkoff.Nominal

	resLastPriceFromTinkoff, err := tinkoffApi.GetLastPriceFromTinkoffInPersentageToNominal(client, s.InstrumentUid)
	if err != nil {
		return errors.New("service:GetSpecificationsFromMoex:" + err.Error())
	}
	// Округляем до двух занков после запятой
	s.SellPrice = math.Round(resLastPriceFromTinkoff/100*s.Nominal*100) / 100
	return nil

}

func ProcessOperations(client *investgo.Client, operations []Operation) (*ReportPositions, error) {
	processPosition := &ReportPositions{}
	for _, operation := range operations {
		switch operation.Type {
		// 2	Удержание НДФЛ по купонам.
		// 8    Удержание налога по дивидендам.
		case WithholdingOfPersonalIncomeTaxOnCoupons, WithholdingOfPersonalIncomeTaxOnDividends:
			err := processWithholdingOfPersonalIncomeTaxOnCouponsOrDividends(operation, processPosition)
			if err != nil {
				return nil, errors.New("ProcessOperations: processWithholdingOfPersonalIncomeTaxOnCouponsOrDividends" + err.Error())
			}

			// 10	Частичное погашение облигаций.
		case PartialRedemptionOfBonds:
			err := processPartialRedemptionOfBonds(operation, processPosition)
			if err != nil {
				return nil, errors.New("ProcessOperations: processPartialRedemptionOfBonds" + err.Error())
			}

			// 15	Покупка ЦБ.
			// 16	Покупка ЦБ с карты.
			// 57   Перевод ценных бумаг с ИИС на Брокерский счет
		case PurchaseOfSecurities, PurchaseOfSecuritiesWithACard, TransferOfSecuritiesFromIISToABrokerageAccount:
			// Проверяем операцию на выполнение.
			// Т.е. операция может быть неисполнена, когда по заявленой цене не было встречного предложения
			if operation.QuantityDone == 0 {
				continue
			} else {
				err := processPurchaseOfSecurities(client, operation, processPosition)
				if err != nil {
					return nil, errors.New("service:ProcessOperations:" + err.Error())
				}
			}
			// 17	Перевод ценных бумаг из другого депозитария.
		case TransferOfSecuritiesFromAnotherDepository:
			// Проверяем операцию на выполнение.
			// Т.е. операция может быть неисполнена, когда по заявленой цене не было встречного предложения
			if operation.QuantityDone == 0 {
				continue
			} else {
				err := processTransferOfSecuritiesFromAnotherDepository(client, operation, processPosition)
				if err != nil {
					return nil, errors.New("service:ProcessOperations:" + err.Error())
				}
			}
			// 19	Удержание комиссии за операцию.
		case WithhouldingACommissionForTheTransaction:
			// Посчитали комисссию в операции покупки(15,16.17,57) и операции продажи(22)

			// 21	Выплата дивидендов.
		case PaymentOfDividends:
			err := processPaymentOfDividends(operation, processPosition)
			if err != nil {
				return nil, errors.New("ProcessOperations: processPaymentOfDividends" + err.Error())
			}
			// 22	Продажа ЦБ.
		case SaleOfSecurities:
			// Проверяем операцию на выполнение.
			// Т.е. операция может быть неисполнена, когда по заявленой цене не было встречного предложения
			if operation.QuantityDone == 0 {
				continue
			} else {
				err := processSellOfSecurities(&operation, processPosition)
				if err != nil {
					return nil, errors.New("ProcessOperations: processSellOfSecurities" + err.Error())
				}
			}

			// 23 Выплата купонов.
		case PaymentOfCoupons:
			err := processPaymentOfCoupons(operation, processPosition)
			if err != nil {
				return nil, errors.New("ProcessOperations: processPaymentOfCoupons" + err.Error())
			}

			// 47	Гербовый сбор.
		case StampDuty:
			err := processStampDuty(operation, processPosition)
			if err != nil {
				return nil, errors.New("ProcessOperations: processStampDuty" + err.Error())
			}
		default:
			continue

		}
	}
	return processPosition, nil
}
