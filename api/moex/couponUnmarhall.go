package moex

import (
	"encoding/json"
	"errors"
)

type Coupons struct {
	Columns []string `json:"columns"`
	Data    []Coupon `json:"data"`
}

type Coupon struct {
	Isin             string   `json:"isin"`
	Name             string   `json:"name"`
	Issuevalue       float64  `json:"issuevalue"`
	Coupondate       string   `json:"coupondate"`
	Recorddate       *string  `json:"recorddate"`
	Startdate        string   `json:"startdate"`
	Initialfacevalue float64  `json:"initialfacevalue"`
	Facevalue        float64  `json:"facevalue"`
	Faceunit         string   `json:"faceunit"`
	Value            *float64 `json:"value,omitempty"`
	Valueprc         *float64 `json:"valueprc,omitempty"`
	Value_rub        *float64 `json:"value_rub,omitempty"`
	Secid            string   `json:"secid"`
	Primary_boardid  string   `json:"primary_boardid"`
}

func (c *Coupon) UnmarshalJSON(data []byte) error {
	dataSlice := make([]any, 14)
	err := json.Unmarshal(data, &dataSlice)
	if err != nil {
		return errors.New("CustomFloat64: UnmarshalJSON: " + err.Error())
	}
	c.Isin = dataSlice[0].(string)
	c.Name = dataSlice[1].(string)
	c.Issuevalue = dataSlice[2].(float64)
	c.Coupondate = dataSlice[3].(string)
	c.Recorddate = checkStringNull(dataSlice[4])
	c.Startdate = dataSlice[5].(string)
	c.Initialfacevalue = dataSlice[6].(float64)
	c.Facevalue = dataSlice[7].(float64)
	c.Faceunit = dataSlice[8].(string)
	c.Value = checkFloa64Null(dataSlice[9])
	c.Valueprc = checkFloa64Null(dataSlice[10])
	c.Value_rub = checkFloa64Null(dataSlice[11])
	c.Secid = dataSlice[12].(string)
	c.Primary_boardid = dataSlice[13].(string)

	return nil
}
