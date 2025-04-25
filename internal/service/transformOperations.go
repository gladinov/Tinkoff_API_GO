package service

import (
	"time"

	"github.com/gothanks/myapp/other_func"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type Operation struct {
	Currency          string
	Cursor            string
	BrokerAccountId   string
	Operation_Id      string
	ParentOperationId string
	Name              string
	Date              time.Time // Время в UTC
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
	QuantityDone      float64
	CancelDateTime    time.Time
	CancelReason      string
	TradesInfo        *pb.OperationItemTrades
	AssetUid          string
	ChildOperations   []*pb.ChildOperationItem
}

type OperationDB struct {
	Name           string
	Date           string
	Figi           string
	Operation_Id   string
	QuantityDone   float64
	InstrumentType string
	InstrumentUid  string
	Price          float64
	Currency       string
	AccruedInt     float64
	Commission     float64
	Payment        float64
}

// Приводим операции к удобной структуре
func TransOperations(operations []*pb.OperationItem) ([]Operation) {
	transformOperations := make([]Operation, 0)
	for _, v := range operations {
		transformOperation := Operation{
			Currency:          v.GetPrice().Currency,
			Cursor:            v.GetCursor(),
			BrokerAccountId:   v.GetBrokerAccountId(),
			Operation_Id:      v.GetId(),
			ParentOperationId: v.GetParentOperationId(),
			Name:              v.GetName(),
			Date:              v.Date.AsTime(),
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
			QuantityDone:      float64(v.GetQuantityDone()),
			CancelDateTime:    v.CancelDateTime.AsTime(),
			CancelReason:      v.GetCancelReason(),
			TradesInfo:        v.GetTradesInfo(),
			AssetUid:          v.GetAssetUid(),
			ChildOperations:   v.GetChildOperations(),
		}

		transformOperations = append(transformOperations, transformOperation)
	}
	return transformOperations
}
