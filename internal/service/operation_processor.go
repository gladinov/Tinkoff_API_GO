package service

import (
	"github.com/gothanks/myapp/api/tinkoff_api"
	"github.com/gothanks/myapp/other_func"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	WithholdingOfPersonalIncomeTaxOnCoupons        = 2   // 2	Удержание НДФЛ по купонам.
	WithholdingOfPersonalIncomeTaxOnDividends      = 8   // 8    Удержание налога по дивидендам.
	PartialRedemptionOfBonds                       = 10  // 10	Частичное погашение облигаций.
	PurchaseOfSecurities                           = 15  // 15	Покупка ЦБ.
	PurchaseOfSecuritiesWithACard                  = 16  // 16	Покупка ЦБ с карты.
	TransferOfSecuritiesFromAnotherDepository      = 17  // 17	Перевод ценных бумаг из другого депозитария.
	WithhouldingACommissionForTheTransaction       = 19  // 19	Удержание комиссии за операцию.
	PaymentOfDividends                             = 21  // 21	Выплата дивидендов.
	SaleOfSecurities                               = 22  // 22	Продажа ЦБ.
	PaymentOfCoupons                               = 23  // 23 Выплата купонов.
	StampDuty                                      = 47  // 47	Гербовый сбор.
	TransferOfSecuritiesFromIISToABrokerageAccount = 57  // 57   Перевод ценных бумаг с ИИС на Брокерский счет
	EuroTransBuyCost                               = 240 //Стоимость Евротранса при переводе из другого депозитария
)

type ReportPositions struct {
	CurrentPositions []SharePosition
	ClosePostion     []SharePosition
}

type SharePosition struct {
	Name               string
	BuyDate            *timestamppb.Timestamp
	SellDate           *timestamppb.Timestamp
	BuyCursor          string
	SellCursor         string
	Quantity           int
	Figi               string
	InstrumentType     string
	InstrumentUid      string
	BuyPrice           float64
	SellPrice          float64
	CurrentPrice       float64
	BuyPayment         float64
	SellPayment        float64
	Currency           string
	AccruedInt         float64 // НКД
	PER                float64 // Частичное досрочное гашение
	TotalCoupon        float64
	TotalDividend      float64
	TotalComission     float64
	TotalTax           float64
	PositionProfit     float64
	ProfitInPercentage float64
}

type Operation struct {
	Currency          string
	Cursor            string
	BrokerAccountId   string
	Operation_Id      string
	ParentOperationId string
	Name              string
	Date              *timestamppb.Timestamp
	Type              int64
	Description       string
	State             int64
	InstrumentUid     string
	Figi              string
	InstrumentType    string
	InstrumentKind    string
	PositionUid       string
	Payment           float64
	Price             float64
	Commission        float64
	Yield             float64
	YieldRelative     float64
	AccruedInt        float64
	Quantity          int64
	QuantityRest      int64
	QuantityDone      int64
	CancelDateTime    *timestamppb.Timestamp
	CancelReason      string
	TradesInfo        *pb.OperationItemTrades
	AssetUid          string
	ChildOperations   []*pb.ChildOperationItem
}

func transOperaions(account *tinkoff_api.Account) []Operation {
	resList := make([]Operation, 0)
	for _, v := range account.Operations {
		transOperationRet := Operation{
			Currency:          v.GetPrice().Currency,
			Cursor:            v.GetCursor(),
			BrokerAccountId:   v.GetBrokerAccountId(),
			Operation_Id:      v.GetId(),
			ParentOperationId: v.GetParentOperationId(),
			Name:              v.GetName(),
			Date:              v.GetDate(),
			Type:              int64(v.GetType()),
			Description:       v.GetDescription(),
			State:             int64(v.GetState()),
			InstrumentUid:     v.GetInstrumentUid(),
			Figi:              v.GetFigi(),
			InstrumentType:    v.GetInstrumentType(),
			InstrumentKind:    string(v.GetInstrumentKind()),
			PositionUid:       v.GetPositionUid(),
			Payment:           other_func.MoneyValue(v.GetPayment()),
			Price:             other_func.MoneyValue(v.GetPrice()),
			Commission:        other_func.MoneyValue(v.GetCommission()),
			Yield:             other_func.MoneyValue(v.GetYield()),
			YieldRelative:     other_func.CastMoney(v.GetYieldRelative()),
			AccruedInt:        other_func.MoneyValue(v.GetAccruedInt()),
			Quantity:          v.GetQuantity(),
			QuantityRest:      v.GetQuantityRest(),
			QuantityDone:      v.GetQuantityDone(),
			CancelDateTime:    v.GetCancelDateTime(),
			CancelReason:      v.GetCancelReason(),
			TradesInfo:        v.GetTradesInfo(),
			AssetUid:          v.GetAssetUid(),
			ChildOperations:   v.GetChildOperations(),
		}
		resList = append(resList, transOperationRet)
	}
	return resList
}

// в данную фунцию можн о дабавить еще текущую цену бумаги и запрашивать ее один раз при полученни UID
func ProcessPositions(operation Operation) (*ReportPositions, error) {
	processPosition := &ReportPositions{}
	countPositions := len(processPosition.CurrentPositions)
	switch operation.Type {
	// 2	Удержание НДФЛ по купонам.
	// 8    Удержание налога по дивидендам.
	case WithholdingOfPersonalIncomeTaxOnCoupons, WithholdingOfPersonalIncomeTaxOnDividends:
		if countPositions != 0 {
			for i := range processPosition.CurrentPositions {
				processPosition.CurrentPositions[i].TotalDividend += operation.Payment / float64(countPositions)
			}
		}
		// 10	Частичное погашение облигаций.
	case PartialRedemptionOfBonds:
		if countPositions != 0 {
			for i := range processPosition.CurrentPositions {
				processPosition.CurrentPositions[i].PER += operation.Payment / float64(countPositions)
			}
		}
		// 15	Покупка ЦБ.
		// 16	Покупка ЦБ с карты.
		// 57   Перевод ценных бумаг с ИИС на Брокерский счет
	case PurchaseOfSecurities, PurchaseOfSecuritiesWithACard, TransferOfSecuritiesFromIISToABrokerageAccount:
		position := SharePosition{
			Name:           operation.Name,
			BuyDate:        operation.Date,
			Figi:           operation.Figi,
			BuyCursor:      operation.Cursor,
			Quantity:       1,
			InstrumentType: operation.InstrumentType,
			InstrumentUid:  operation.InstrumentUid,
			BuyPrice:       operation.Price,
			Currency:       operation.Currency,
		}
		// НКД при покупке добавляем с отрицательным знаком, т.к. потом получим эту сумму с купоном или при продаже
		if operation.QuantityDone != 0 {
			position.AccruedInt = -operation.AccruedInt / float64(operation.QuantityDone)
		}
		// Проверяем на ноль значение комиссии и добавляем
		if operation.QuantityDone != 0 {
			position.TotalComission = operation.Commission / float64(operation.QuantityDone)
		}
		// Для каждой купленной бумаги добавляю SharePosition
		for i := int64(0); i < operation.QuantityDone; i++ {
			processPosition.CurrentPositions = append(processPosition.CurrentPositions, position)
		}
		// 17	Перевод ценных бумаг из другого депозитария.
	case TransferOfSecuritiesFromAnotherDepository:
		position := SharePosition{
			Name:           operation.Name,
			BuyDate:        operation.Date,
			BuyCursor:      operation.Cursor,
			Quantity:       1,
			Figi:           operation.Figi,
			InstrumentType: operation.InstrumentType,
			InstrumentUid:  operation.InstrumentUid,
			BuyPrice:       operation.Price,
			Currency:       operation.Currency,
		}
		// НКД при покупке добавляем с отрицательным знаком, т.к. потом получим эту сумму с купоном или при продаже
		if operation.QuantityDone != 0 {
			position.AccruedInt = -operation.AccruedInt / float64(operation.QuantityDone)
		}
		// Проверяем на ноль значение комиссии и добавляем
		if operation.QuantityDone != 0 {
			position.TotalComission = operation.Commission / float64(operation.QuantityDone)
		}
		// Для Евротранса исключение
		if operation.InstrumentUid == "02b2ea14-3c4b-47e8-9548-45a8dbcc8f8a" {
			position.BuyPrice = EuroTransBuyCost
		}
		// Для каждой купленной бумаги добавляю SharePosition
		for i := int64(0); i < operation.QuantityDone; i++ {
			processPosition.CurrentPositions = append(processPosition.CurrentPositions, position)
		}
		// 19	Удержание комиссии за операцию.
	case WithhouldingACommissionForTheTransaction:
		// Посчитали комисссию в операции покупки(15) и операции продажи(22)

		// 21	Выплата дивидендов.
	case PaymentOfDividends:
		if countPositions != 0 {
			for i := range processPosition.CurrentPositions {
				processPosition.CurrentPositions[i].TotalDividend += operation.Payment / float64(countPositions)
			}
		}
		// 22	Продажа ЦБ.
	case SaleOfSecurities:
		quantitySell := operation.QuantityDone
		for i := range int(quantitySell) {
			// Сократим название переменной
			currPostion := &processPosition.CurrentPositions[i]
			// Устанавливаем НКД
			if operation.QuantityDone != 0 {
				currPostion.AccruedInt += operation.AccruedInt / float64(operation.QuantityDone)
			}
			// Устанавливаем цену продажи для бумаги
			currPostion.SellPrice = operation.Price
			// Устанавливаем курсор продажи
			currPostion.SellCursor = operation.Cursor
			// Рассчитываем срок владения
			currPostion.SellDate = operation.Date
			buyDate := currPostion.BuyDate.AsTime()
			sellDate := currPostion.SellDate.AsTime()
			timeDuration := sellDate.Sub(buyDate)
			threeYearHour := 26304
			// Плюсуем комиссию за продажу бумаг
			currPostion.TotalComission += operation.Commission / float64(operation.QuantityDone)
			// Рассчитываем налог с продажи бумаги, если сумма продажи больше суммы покупки
			sellPrice := currPostion.SellPrice
			buyPrice := currPostion.BuyPrice
			Profit := buyPrice - sellPrice
			if Profit > 0 && timeDuration.Hours() > float64(threeYearHour) {
				currPostion.TotalTax += Profit * 0.13
			}
			// Считаем Доход
			currPostion.PositionProfit = currPostion.SellPrice + currPostion.TotalDividend - currPostion.TotalComission - currPostion.BuyPrice - currPostion.TotalTax + currPostion.AccruedInt + currPostion.PER + currPostion.TotalCoupon
			currPostion.ProfitInPercentage = currPostion.PositionProfit / currPostion.BuyPrice

		}
		processPosition.ClosePostion = append(processPosition.ClosePostion, processPosition.CurrentPositions[:quantitySell]...)
		processPosition.CurrentPositions = processPosition.CurrentPositions[quantitySell:]
		// 23 Выплата купонов.
	case PaymentOfCoupons:
		if countPositions != 0 {
			for i := range processPosition.CurrentPositions {
				processPosition.CurrentPositions[i].TotalCoupon += operation.Payment / float64(countPositions)
			}
		}
		// 47	Гербовый сбор.
	case StampDuty:
		if countPositions != 0 {
			for i := range processPosition.CurrentPositions {
				processPosition.CurrentPositions[i].TotalComission += operation.Payment / float64(countPositions)
			}
		}
	}
	return processPosition, nil
}

func UnionPositions(positions []SharePosition) []SharePosition {
	if len(positions) == 0 {
		return nil
	}
	buyCursor := positions[0].BuyCursor
	sellCursor := positions[0].SellCursor
	retPosition := positions[0]
	resList := make([]SharePosition, 0)
	for i := 1; i < len(positions); i++ {
		if positions[i].BuyCursor == buyCursor && positions[i].SellCursor == sellCursor {
			retPosition.BuyPayment += positions[i].BuyPrice
			retPosition.SellPayment += positions[i].SellPrice
			retPosition.TotalDividend += positions[i].TotalDividend
			retPosition.TotalComission += positions[i].TotalComission
			retPosition.TotalTax += positions[i].TotalTax
			retPosition.Quantity += positions[i].Quantity
			retPosition.PositionProfit += positions[i].PositionProfit
			if len(positions)-1 == i {
				resList = append(resList, retPosition)
			}
		} else {
			resList = append(resList, retPosition)
			buyCursor = positions[i].BuyCursor
			sellCursor = positions[i].SellCursor
			retPosition = positions[i]
			if len(positions)-1 == i {
				resList = append(resList, retPosition)
			}
		}
	}
	return resList
}
