package moex

import (
	"encoding/json"
	"errors"
	"time"
)

type Durations struct {
	Marketdata_yields Marketdata_yields
}

type Marketdata_yields struct {
	Data Duration
}

type Duration struct {
	Duration float64
	Date     time.Time
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	dataSlice := make([][]any, 0)
	err := json.Unmarshal(data, &dataSlice)
	if err != nil {
		return errors.New("CustomFloat64: UnmarshalJSON: " + err.Error())
	}
	d.Date = time.Now()
	d.Duration = dataSlice[0][0].(float64)
	return nil
}
