package service

import "errors"

// 2	Удержание НДФЛ по купонам.
// 8    Удержание налога по дивидендам.
func processWithholdingOfPersonalIncomeTaxOnCouponsOrDividends(operation Operation, processPosition *ReportPositions) error {
	if processPosition.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for _, currentPosition := range processPosition.CurrentPositions {
			proportion := currentPosition.Quantity / processPosition.Quantity
			currentPosition.TotalTax += operation.Payment * proportion
		}
	}
	return nil
}

// 10	Частичное погашение облигаций.
func processPartialRedemptionOfBonds(operation Operation, processPosition *ReportPositions) error {
	if processPosition.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for _, currentPosition := range processPosition.CurrentPositions {
			proportion := currentPosition.Quantity / processPosition.Quantity
			currentPosition.PER += operation.Payment * proportion
		}
	}
	return nil

}

// 15	Покупка ЦБ.
// 16	Покупка ЦБ с карты.
// 57   Перевод ценных бумаг с ИИС на Брокерский счет
func processPurchaseOfSecurities(operation Operation, processPosition *ReportPositions) {
	// при обработке фьючерсов и акций, где была маржтнальная позиция,
	//  функцию надо переделать так, чтобы проверялось наличие позиций с отрицательным количеством бумаг(коротких позиций)
	position := SharePosition{
		Name:           operation.Name,
		BuyDate:        operation.Date,
		Figi:           operation.Figi,
		BuyOperationID: operation.Operation_Id,
		Quantity:       operation.QuantityDone,
		InstrumentType: operation.InstrumentType,
		InstrumentUid:  operation.InstrumentUid,
		BuyPrice:       operation.Price,
		Currency:       operation.Currency,
		BuyAccruedInt:  operation.AccruedInt, // НКД при покупке
		TotalComission: operation.Commission,
	}
	processPosition.CurrentPositions = append(processPosition.CurrentPositions, position)
	processPosition.Quantity += operation.QuantityDone
}

// 17	Перевод ценных бумаг из другого депозитария.
func processTransferOfSecuritiesFromAnotherDepository(operation Operation, processPosition *ReportPositions) {
	// при обработке фьючерсов и акций, где была маржтнальная позиция,
	//  функцию надо переделать так, чтобы проверялось наличие позиций с отрицательным количеством бумаг(коротких позиций)
	position := SharePosition{
		Name:           operation.Name,
		BuyDate:        operation.Date,
		Figi:           operation.Figi,
		BuyOperationID: operation.Operation_Id,
		Quantity:       operation.QuantityDone,
		InstrumentType: operation.InstrumentType,
		InstrumentUid:  operation.InstrumentUid,
		BuyPrice:       operation.Price,
		Currency:       operation.Currency,
		BuyAccruedInt:  operation.AccruedInt, // НКД при покупке
		TotalComission: operation.Commission,
	}
	// Для Евротранса исключение
	if operation.InstrumentUid == "02b2ea14-3c4b-47e8-9548-45a8dbcc8f8a" {
		position.BuyPrice = EuroTransBuyCost
	}
	processPosition.CurrentPositions = append(processPosition.CurrentPositions, position)
	processPosition.Quantity += operation.QuantityDone
}

// 21	Выплата дивидендов.
func processPaymentOfDividends(operation Operation, processPosition *ReportPositions) error {
	if processPosition.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range processPosition.CurrentPositions {
			currentPosition := &processPosition.CurrentPositions[i]
			proportion := currentPosition.Quantity / processPosition.Quantity
			// Минус, т.к. operation.Payment с отрицательным знаком в отчете
			currentPosition.TotalDividend += operation.Payment * proportion
		}
	}
	return nil
}

// 23 Выплата купонов.
func processPaymentOfCoupons(operation Operation, processPosition *ReportPositions) error {
	if processPosition.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range processPosition.CurrentPositions {
			currentPosition := &processPosition.CurrentPositions[i]
			proportion := currentPosition.Quantity / processPosition.Quantity
			// Минус, т.к. operation.Payment с отрицательным знаком в отчете
			currentPosition.TotalCoupon += operation.Payment * proportion
		}
	}
	return nil
}

// 47	Гербовый сбор.
func processStampDuty(operation Operation, processPosition *ReportPositions) error {
	if processPosition.Quantity == 0 {
		return errors.New("divide by zero")
	} else {
		for i := range processPosition.CurrentPositions {
			currentPosition := &processPosition.CurrentPositions[i]
			proportion := currentPosition.Quantity / processPosition.Quantity
			// Минус, т.к. operation.Payment с отрицательным знаком в отчете
			currentPosition.TotalComission += operation.Payment * proportion
		}
	}
	return nil
}
