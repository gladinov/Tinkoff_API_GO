package other_func

import (
	"errors"
	"time"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CastMoney(v *pb.Quotation) float64 {
	if v != nil {
		r := float64(v.Units) + float64(v.Nano/1e9)
		return r
	}
	return 0
}

func MoneyValue(v *pb.MoneyValue) float64 {
	if v != nil {
		r := float64(v.Units) + float64(v.Nano/1e9)
		return r
	}
	return 0
}

func ProtoToTime(ts *timestamppb.Timestamp) (time.Time, error) {
	if err := ts.CheckValid(); err != nil {
		return time.Time{}, errors.New("invalid timestamp" + err.Error())
	}
	return ts.AsTime(), nil
}

func StringRFC3339ToTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, errors.New("StringRFC3339ToTime: time.Parse" + err.Error())
	}
	return t, nil
}
