package moex

import (
	"encoding/json"
	"errors"
)

type Offers struct {
	Columns []string `json:"columns"`
	Data    []Offer  `json:"data"`
}

type Offer struct {
	Isin            string   `json:"isin"`
	Name            string   `json:"name"`
	Issuevalue      float64  `json:"issuevalue"`
	Offerdate       string   `json:"offerdate"`
	Offerdatestart  string   `json:"offerdatestart"`
	Offerdateend    string   `json:"offerdateend"`
	Facevalue       float64  `json:"facevalue"`
	Faceunit        string   `json:"faceunit"`
	Price           *float64 `json:"price"`
	Value           float64  `json:"value"`
	Agent           *string  `json:"agent,omitempty"`
	Offertype       string   `json:"offertype"`
	Secid           string   `json:"secid"`
	Primary_boardid string   `json:"primary_boardid"`
}

func (o *Offer) UnmarshalJSON(data []byte) error {
	dataSlice := make([]any, 14) // Количество элементов соответствует полям структуры
	err := json.Unmarshal(data, &dataSlice)
	if err != nil {
		return errors.New("Offer: UnmarshalJSON: " + err.Error())
	}

	o.Isin = dataSlice[0].(string)
	o.Name = dataSlice[1].(string)
	o.Issuevalue = dataSlice[2].(float64)
	o.Offerdate = dataSlice[3].(string)
	o.Offerdatestart = dataSlice[4].(string)
	o.Offerdateend = dataSlice[5].(string)
	o.Facevalue = dataSlice[6].(float64)
	o.Faceunit = dataSlice[7].(string)
	o.Price = checkFloa64Null(dataSlice[8])
	o.Value = dataSlice[9].(float64)
	o.Agent = checkStringNull(dataSlice[10])
	o.Offertype = dataSlice[11].(string)
	o.Secid = dataSlice[12].(string)
	o.Primary_boardid = dataSlice[13].(string)

	return nil
}
