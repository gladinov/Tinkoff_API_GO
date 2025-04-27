package moex

import (
	"encoding/json"
	"errors"
)

type Amortizations struct {
	Columns []string       `json:"columns"`
	Data    []Amortization `json:"data"`
}

type Amortization struct {
	Isin             string  `json:"isin"`
	Name             string  `json:"name"`
	Issuevalue       float64 `json:"issuevalue"`
	Amortdate        string  `json:"amortdate"`
	Facevalue        float64 `json:"facevalue"`
	Initialfactvalue float64 `json:"initialfactvalue"`
	Faceuint         string  `json:"faceunit"`
	Valueprc         float64 `json:"valueprc"`
	Value            float64 `json:"value"`
	Value_rub        float64 `json:"value_rub"`
	Data_source      string  `json:"data_source"`
	Secid            string  `json:"secid"`
	Primary_boardid  string  `json:"primary_boardid"`
}

func (a *Amortization) UnmarshalJSON(data []byte) error {
	dataSlice := make([]any, 13)
	err := json.Unmarshal(data, &dataSlice)
	if err != nil {
		return errors.New("CustomFloat64: UnmarshalJSON: " + err.Error())
	}
	a.Isin = dataSlice[0].(string)
	a.Name = dataSlice[1].(string)
	a.Issuevalue = dataSlice[2].(float64)
	a.Amortdate = dataSlice[3].(string)
	a.Facevalue = dataSlice[4].(float64)
	a.Initialfactvalue = dataSlice[5].(float64)
	a.Faceuint = dataSlice[6].(string)
	a.Valueprc = dataSlice[7].(float64)
	a.Value = dataSlice[8].(float64)
	a.Value_rub = dataSlice[9].(float64)
	a.Data_source = dataSlice[10].(string)
	a.Secid = dataSlice[11].(string)
	a.Primary_boardid = dataSlice[12].(string)

	return nil

}
