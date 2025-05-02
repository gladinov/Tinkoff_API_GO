package service

import (
	"time"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type Operation struct {
	Currency          string
	BrokerAccountId   string
	Operation_Id      string
	ParentOperationId string
	Name              string
	Date              time.Time // Время в UTC
	Type              int64
	Description       string
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
	AssetUid          string
}

// Приводим операции к удобной структуре
func TransOperations(operations []*pb.OperationItem) []Operation {
	transformOperations := make([]Operation, 0)
	for _, v := range operations {
		transformOperation := Operation{
			Currency:          v.GetPrice().Currency,
			BrokerAccountId:   v.GetBrokerAccountId(),
			Operation_Id:      v.GetId(),
			ParentOperationId: v.GetParentOperationId(),
			Name:              v.GetName(),
			Date:              v.Date.AsTime(),
			Type:              int64(v.GetType()),
			Description:       v.GetDescription(),
			InstrumentUid:     v.GetInstrumentUid(),
			Figi:              v.GetFigi(),
			InstrumentType:    v.GetInstrumentType(),
			InstrumentKind:    string(v.GetInstrumentKind()),
			PositionUid:       v.GetPositionUid(),
			Payment:           v.GetPayment().ToFloat(),
			Price:             v.GetPrice().ToFloat(),
			Commission:        v.GetCommission().ToFloat(),
			Yield:             v.GetYield().ToFloat(),
			YieldRelative:     v.GetYieldRelative().ToFloat(),
			AccruedInt:        v.GetAccruedInt().ToFloat(),
			QuantityDone:      float64(v.GetQuantityDone()),

			AssetUid: v.GetAssetUid(),
		}
		transformOperations = append(transformOperations, transformOperation)
	}
	return transformOperations
}
