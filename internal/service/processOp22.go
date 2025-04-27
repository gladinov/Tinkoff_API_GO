package service

import (
	"errors"
)

// 22	Продажа ЦБ.
func processSellOfSecurities(operation *Operation, processPosition *ReportPositions) error {
	processPosition.Quantity -= operation.QuantityDone
	// Переписать ПОЗЖЕ Переменная deleteCount отслеживает кол-во закрытых позиций для дальнейшего удаления в которых Кол-во проданных
	// бумаг больше кол-ва бумаг в текущей позиции
	var deleteCount int
end:
	for i := range processPosition.CurrentPositions {
		currPosition := &processPosition.CurrentPositions[i]
		currentQuantity := currPosition.Quantity
		sellQuantity := operation.QuantityDone
		switch {
		case currentQuantity > sellQuantity:
			err := isCurrentQuantityGreaterThanSellQuantity(operation, currPosition, processPosition)
			if err != nil {
				return errors.New("service.isCurrentQuantityGreaterThanSellQuantity" + err.Error())
			}
			// Прерываем цикл
			break end
		case currPosition.Quantity == operation.QuantityDone:
			err := isEqualCurrentQuantityAndSellQuantity(operation, currPosition, processPosition)
			if err != nil {
				return errors.New("service.isEqualCurrentQuantityAndSellQuantity" + err.Error())
			}
			// Прерываем цикл
			break end
		case currentQuantity < sellQuantity:
			// Переменная deleteCount отслеживает кол-во закрытых позиций для дальнейшего удаления
			deleteCount += 1
			err := isCurrentQuantityLessThanSellQuantity(operation, currPosition, processPosition)
			if err != nil {
				return errors.New("service.isCurrentQuantityLessThanSellQuantity" + err.Error())
			}
		}

	}
	// удаляем закрытые позиции из среза текущих позиций
	processPosition.CurrentPositions = processPosition.CurrentPositions[deleteCount:]
	return nil
}

func isCurrentQuantityGreaterThanSellQuantity(operation *Operation, currPosition *SharePosition, processPosition *ReportPositions) error {
	currentQuantity := currPosition.Quantity
	sellQuantity := operation.QuantityDone
	var proportion float64
	if currentQuantity != 0 {
		proportion = sellQuantity / currentQuantity
	} else {
		return errors.New("divide by zero")
	}
	// Создаем закрытую позицию
	closePosition := createClosePosition(*currPosition, *operation)
	// Устанавливаем закрытой позиции кол-во бумаг, которые проданы
	closePosition.Quantity = sellQuantity
	// НКД продажи (На пропорцию не умножаем, т.к. это отдельное поле)
	closePosition.SellAccruedInt = operation.AccruedInt
	// НКД покупки, пересчитанное по пропорции
	closePosition.BuyAccruedInt *= proportion
	// Плюсуем комиссию за продажу бумаг
	closePosition.TotalComission = closePosition.TotalComission*proportion + operation.Commission

	// Считаем Доход без налога
	positionProfit := getSecurityIncomeWithoutTax(closePosition)

	// Считаем налог
	totalTax := getTotalTaxFromPosition(closePosition, positionProfit)

	// Считаем доход после налогообложения
	profitAfterTax := getSecurityIncome(positionProfit, totalTax)

	// Заполняем поля структуры
	closePosition.TotalTax = totalTax
	closePosition.PositionProfit = profitAfterTax

	// Считаем процентный доход
	profitInPercentage, err := getProfitInPercentage(closePosition)
	if err != nil {
		return errors.New("service.getProfitInPercentage" + err.Error())
	}
	closePosition.ProfitInPercentage = profitInPercentage

	// Добавляем закрытую позицию в срез закрытых позиций
	processPosition.ClosePositions = append(processPosition.ClosePositions, closePosition)
	// Отнимаем кол-во проданных бумаг от количества бумаг в текущей позиции
	currPosition.Quantity -= sellQuantity
	// Изменяем значения текущей позиции, умножая на остаток от пропорции
	currPosition.TotalComission = currPosition.TotalComission * (1 - proportion)
	currPosition.PaidTax = currPosition.PaidTax * (1 - proportion)
	currPosition.BuyAccruedInt = currPosition.BuyAccruedInt * (1 - proportion)
	return nil
}

func isEqualCurrentQuantityAndSellQuantity(operation *Operation, currPosition *SharePosition, processPosition *ReportPositions) error {
	// Создаем закрытую позицию
	closePosition := createClosePosition(*currPosition, *operation)
	// НКД
	closePosition.SellAccruedInt = operation.AccruedInt
	// Плюсуем комиссию за продажу бумаг
	closePosition.TotalComission = closePosition.TotalComission + operation.Commission
	// Считаем Доход без налога
	positionProfit := getSecurityIncomeWithoutTax(closePosition)

	// Считаем налог
	totalTax := getTotalTaxFromPosition(closePosition, positionProfit)

	// Считаем доход после налогообложения
	profitAfterTax := getSecurityIncome(positionProfit, totalTax)

	// Заполняем поля структуры
	closePosition.TotalTax = totalTax
	closePosition.PositionProfit = profitAfterTax

	// Считаем процентный доход
	profitInPercentage, err := getProfitInPercentage(closePosition)
	if err != nil {
		return errors.New("service.getProfitInPercentage" + err.Error())
	}
	closePosition.ProfitInPercentage = profitInPercentage

	// Добавляем закрытую позицию в срез закрытых позиций
	processPosition.ClosePositions = append(processPosition.ClosePositions, closePosition)
	processPosition.CurrentPositions = processPosition.CurrentPositions[1:]
	return nil
}

func isCurrentQuantityLessThanSellQuantity(operation *Operation, currPosition *SharePosition, processPosition *ReportPositions) error {
	currentQuantity := currPosition.Quantity
	sellQuantity := operation.QuantityDone
	var proportion float64
	if sellQuantity != 0 {
		proportion = currentQuantity / sellQuantity
	} else {
		return errors.New("divide by zero")
	}
	// Создаем операцию продажи бумаги
	// Создаем закрытую позицию
	closePosition := createClosePosition(*currPosition, *operation)
	// НКД
	closePosition.SellAccruedInt = operation.AccruedInt * proportion
	operation.AccruedInt -= operation.AccruedInt * proportion
	// Плюсуем комиссию за продажу бумаг
	closePosition.TotalComission = closePosition.TotalComission + (operation.Commission * proportion)
	operation.Commission -= operation.Commission * proportion

	// Считаем Доход без налога
	positionProfit := getSecurityIncomeWithoutTax(closePosition)

	// Считаем налог
	totalTax := getTotalTaxFromPosition(closePosition, positionProfit)

	// Считаем доход после налогообложения
	profitAfterTax := getSecurityIncome(positionProfit, totalTax)

	// Заполняем поля структуры
	closePosition.TotalTax = totalTax
	closePosition.PositionProfit = profitAfterTax

	// Считаем процентный доход
	profitInPercentage, err := getProfitInPercentage(closePosition)
	if err != nil {
		return errors.New("service.getProfitInPercentage" + err.Error())
	}
	closePosition.ProfitInPercentage = profitInPercentage

	// Добавляем закрытую позицию в срез закрытых позиций
	processPosition.ClosePositions = append(processPosition.ClosePositions, closePosition)
	// Изменяем значение Quantity.Operation
	operation.QuantityDone -= currPosition.Quantity

	return nil
}

// Доход по позиции до вычета налога
func getSecurityIncomeWithoutTax(p SharePosition) float64 {
	quantity := p.Quantity
	buySellDifference := (p.SellPrice-p.BuyPrice)*quantity + p.SellAccruedInt - p.BuyAccruedInt
	cashFlow := p.TotalCoupon + p.TotalDividend
	positionProfit := buySellDifference + cashFlow + p.TotalComission + p.PartialEarlyRepayment
	return positionProfit
}

// Расход полного налога по закрытой позиции
func getTotalTaxFromPosition(p SharePosition, profit float64) float64 {
	// Рассчитываем срок владения
	buyDate := p.BuyDate
	sellDate := p.SellDate
	timeDuration := sellDate.Sub(buyDate).Hours()
	var tax float64
	// Рассчитываем налог с продажи бумаги, если сумма продажи больше суммы покупки
	if profit > 0 && timeDuration < float64(threeYearInHours) {
		tax = profit * baseTaxRate
	} else {
		tax = 0
	}
	return tax
}

// Расчет прибыли после налогообложения
func getSecurityIncome(profit, tax float64) float64 {
	profitAfterTax := profit - tax
	return profitAfterTax
}

func getProfitInPercentage(p SharePosition) (float64, error) {
	if p.BuyPrice != 0 || p.Quantity != 0 {
		profitInPercentage := p.PositionProfit / (p.BuyPrice * p.Quantity)
		return profitInPercentage, nil
	} else {
		return 0, errors.New("divide by zero")
	}
}

// Создаем закрытую позицию
func createClosePosition(currentPosition SharePosition, operation Operation) SharePosition {
	closePosition := currentPosition
	closePosition.SellPrice = operation.Price
	closePosition.SellPayment = operation.Payment
	closePosition.SellOperationID = operation.Operation_Id
	closePosition.SellDate = operation.Date
	return closePosition
}
