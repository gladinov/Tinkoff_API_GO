package moex

import (
	"encoding/json"
	"errors"
)

type Yields struct {
	History *History `json:"history"`
}

type History struct {
	Data []Values `json:"data"`
}

type Values struct {
	TradeDate       *string  `json:"TRADEDATE"`    // Торговая дата(на момент которой рассчитаны остальные данные)
	MaturityDate    *string  `json:"MATDATE"`      // Дата погашения
	OfferDate       *string  `json:"OFFERDATE"`    // Дата Оферты
	BuybackDate     *string  `json:"BUYBACKDATE"`  // дата обратного выкупа
	YieldToMaturity *float64 `json:"YIELDCLOSE"`   // Доходность к погашению при покупке
	YieldToOffer    *float64 `json:"YIELDTOOFFER"` // Доходность к оферте при покупке
	FaceValue       *float64 `json:"FACEVALUE"`    // номинальная стоимость облигации
	Duration        *float64 `json:"DURATION"`     // дюрация (средневзвешенный срок платежей)

}

func (d *Values) UnmarshalJSON(data []byte) error {
	dataSlice := make([]any, 8)
	err := json.Unmarshal(data, &dataSlice)
	if err != nil {
		return errors.New("CustomFloat64: UnmarshalJSON: " + err.Error())
	}
	d.TradeDate = checkStringNull(dataSlice[0])
	d.MaturityDate = checkStringNull(dataSlice[1])
	d.OfferDate = checkStringNull(dataSlice[2])
	d.BuybackDate = checkStringNull(dataSlice[3])
	d.YieldToMaturity = checkFloa64Null(dataSlice[4])
	d.YieldToOffer = checkFloa64Null(dataSlice[5])
	d.FaceValue = checkFloa64Null(dataSlice[6])
	d.Duration = checkFloa64Null(dataSlice[7])

	return nil
}
