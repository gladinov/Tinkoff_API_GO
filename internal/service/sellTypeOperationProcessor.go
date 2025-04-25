package service

import "errors"

// 22	Продажа ЦБ.
func processSellOfSecurities(operation Operation, processPosition *ReportPositions) error {
	processPosition.Quantity -= operation.QuantityDone
end:
	for i := range processPosition.CurrentPositions {
		currPosition := &processPosition.CurrentPositions[i]
		buyQuantity := currPosition.Quantity
		sellQuantity := operation.QuantityDone
		switch {
		case buyQuantity > sellQuantity:
			currPosition.Quantity -= sellQuantity
			// Создаем операцию продажи бумаги
			closePosition := currPosition
			closePosition.Quantity = operation.QuantityDone
			closePosition.SellPrice = operation.Price
			closePosition.SellPayment = operation.Payment
			closePosition.SellOperationID = operation.Operation_Id
			// НКД
			if buyQuantity != 0 {
				closePosition.SellAccruedInt = operation.AccruedInt
			} else {
				return errors.New("divide by zero")
			}
			// Рассчитываем срок владения
			closePosition.SellDate = operation.Date
			buyDate := closePosition.BuyDate
			sellDate := closePosition.SellDate
			timeDuration := sellDate.Sub(buyDate)
			// Плюсуем комиссию за продажу бумаг
			if buyQuantity != 0 {
				closePosition.TotalComission = closePosition.TotalComission*float64(sellQuantity)/float64(buyQuantity) + operation.Commission
			} else {
				return errors.New("divide by zero")
			}
			// Рассчитываем налог с продажи бумаги, если сумма продажи больше суммы покупки
			sellPrice := closePosition.SellPrice
			buyPrice := closePosition.BuyPrice
			Profit := buyPrice - sellPrice
			if Profit > 0 && timeDuration.Hours() < float64(threeYearInHour) {
				closePosition.TotalTax -= (Profit * float64(closePosition.Quantity) * 0.13)
			}
			// Считаем Доход
			closePosition.PositionProfit = float64(sellQuantity)*(closePosition.SellPrice-closePosition.BuyPrice) + closePosition.TotalDividend + closePosition.TotalComission + closePosition.TotalTax + closePosition.SellAccruedInt - closePosition.BuyAccruedInt
			if closePosition.BuyPrice != 0 {
				closePosition.ProfitInPercentage = closePosition.PositionProfit / closePosition.BuyPrice
			} else {
				return errors.New("divide by zero")
			}
			// Добавляем закрытую позицию в срез закрытых позиций
			processPosition.ClosePositions = append(processPosition.ClosePositions, *closePosition)
			// Прерываем цикл
			break end
		case currPosition.Quantity == operation.QuantityDone:
			// Создаем операцию продажи бумаги
			closePosition := currPosition
			closePosition.SellPrice = operation.Price
			closePosition.SellPayment = operation.Payment
			closePosition.SellOperationID = operation.Operation_Id
			// НКД
			closePosition.SellAccruedInt = operation.AccruedInt
			// Плюсуем комиссию за продажу бумаг
			closePosition.TotalComission = closePosition.TotalComission + operation.Commission
			// Рассчитываем срок владения
			closePosition.SellDate = operation.Date
			buyDate := closePosition.BuyDate
			sellDate := closePosition.SellDate
			timeDuration := sellDate.Sub(buyDate)
			// Рассчитываем налог с продажи бумаги, если сумма продажи больше суммы покупки
			sellPrice := closePosition.SellPrice
			buyPrice := closePosition.BuyPrice
			Profit := buyPrice - sellPrice
			if Profit > 0 && timeDuration.Hours() < float64(threeYearInHour) {
				closePosition.TotalTax -= (Profit * float64(closePosition.Quantity) * 0.13)
			}
			// Считаем Доход
			closePosition.PositionProfit = float64(sellQuantity)*(closePosition.SellPrice-closePosition.BuyPrice) + closePosition.TotalDividend + closePosition.TotalComission + closePosition.TotalTax + closePosition.SellAccruedInt - closePosition.BuyAccruedInt
			if closePosition.BuyPrice != 0 {
				closePosition.ProfitInPercentage = closePosition.PositionProfit / closePosition.BuyPrice
			} else {
				return errors.New("divide by zero")
			}
			// Добавляем закрытую позицию в срез закрытых позиций
			processPosition.ClosePositions = append(processPosition.ClosePositions, *closePosition)
			processPosition.CurrentPositions = processPosition.CurrentPositions[1:]
			// Прерываем цикл
			break end
		case currPosition.Quantity < operation.QuantityDone:
			proportion := currPosition.Quantity / operation.QuantityDone
			// Создаем операцию продажи бумаги
			closePosition := currPosition
			closePosition.SellPrice = operation.Price
			closePosition.SellPayment = operation.Payment
			closePosition.SellOperationID = operation.Operation_Id
			// НКД
			closePosition.SellAccruedInt = operation.AccruedInt * proportion
			operation.AccruedInt -= operation.AccruedInt * proportion
			// Плюсуем комиссию за продажу бумаг
			closePosition.TotalComission = closePosition.TotalComission + (operation.Commission * proportion)
			operation.Commission -= operation.Commission * proportion
			// Рассчитываем срок владения
			closePosition.SellDate = operation.Date
			buyDate := closePosition.BuyDate
			sellDate := closePosition.SellDate
			timeDuration := sellDate.Sub(buyDate)
			// Рассчитываем налог с продажи бумаги, если сумма продажи больше суммы покупки
			sellPrice := closePosition.SellPrice
			buyPrice := closePosition.BuyPrice
			Profit := buyPrice - sellPrice
			if Profit > 0 && timeDuration.Hours() < float64(threeYearInHour) {
				closePosition.TotalTax -= (Profit * float64(closePosition.Quantity) * 0.13)
			}
			// Считаем Доход
			closePosition.PositionProfit = float64(sellQuantity)*(closePosition.SellPrice-closePosition.BuyPrice) + closePosition.TotalDividend + closePosition.TotalComission + closePosition.TotalTax + closePosition.SellAccruedInt - closePosition.BuyAccruedInt
			if closePosition.BuyPrice != 0 {
				closePosition.ProfitInPercentage = closePosition.PositionProfit / closePosition.BuyPrice
			} else {
				return errors.New("divide by zero")
			}
			// Добавляем закрытую позицию в срез закрытых позиций
			processPosition.ClosePositions = append(processPosition.ClosePositions, *closePosition)
			// Изменяем значение Quantity.Operation
			operation.QuantityDone = currPosition.Quantity
		}
	}
	return nil
}
