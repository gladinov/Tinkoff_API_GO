package main

import (
	"fmt"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Operation struct {
	Currency          string
	Cursor            string
	BrokerAccountId   string
	Operation_Id      string
	ParentOperationId string
	Name              string
	Date              *timestamppb.Timestamp
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
	Quantity          int64
	QuantityRest      int64
	QuantityDone      int64
	CancelDateTime    *timestamppb.Timestamp
	CancelReason      string
	TradesInfo        *pb.OperationItemTrades
	AssetUid          string
	ChildOperations   []*pb.ChildOperationItem
}

func GetOpp(
	logger *zap.SugaredLogger,
	operationsService *investgo.OperationsServiceClient,
	account *Account) {

	operationsResp, err := operationsService.GetOperationsByCursor(&investgo.GetOperationsByCursorRequest{
		AccountId: account.Id,
		Limit:     1000,
	})
	if err != nil {
		logger.Errorf(err.Error())
	} else {
		ops := operationsResp.GetOperationsByCursorResponse.Items
		transOperaions(ops, account)
		nextCursor := operationsResp.NextCursor
		for nextCursor != "" {
			operationsResp, err := operationsService.GetOperationsByCursor(&investgo.GetOperationsByCursorRequest{
				AccountId: account.Id,
				Limit:     1000,
				Cursor:    nextCursor,
			})
			if err != nil {
				logger.Errorf(err.Error())
			} else {
				nextCursor = operationsResp.NextCursor
				ops := operationsResp.GetOperationsByCursorResponse.Items
				transOperaions(ops, account)

			}
		}

	}
	// fmt.Println(account.Operations)
	fmt.Printf("✓ Добавлено %v операции в Account.Operation по счету %s\n", len(account.Operations), account.Id)
}

func transOperaions(operationT []*pb.OperationItem, account *Account) {
	for _, v := range operationT {
		transOperationRet := Operation{
			Currency:          v.GetPrice().Currency,
			Cursor:            v.GetCursor(),
			BrokerAccountId:   v.GetBrokerAccountId(),
			Operation_Id:      v.GetId(),
			ParentOperationId: v.GetParentOperationId(),
			Name:              v.GetName(),
			Date:              v.GetDate(),
			Type:              int64(v.GetType()),
			Description:       v.GetDescription(),
			State:             int64(v.GetState()),
			InstrumentUid:     v.GetInstrumentUid(),
			Figi:              v.GetFigi(),
			InstrumentType:    v.GetInstrumentType(),
			InstrumentKind:    string(v.GetInstrumentKind()),
			PositionUid:       v.GetPositionUid(),
			Payment:           MoneyValue(v.GetPayment()),
			Price:             MoneyValue(v.GetPrice()),
			Commission:        MoneyValue(v.GetCommission()),
			Yield:             MoneyValue(v.GetYield()),
			YieldRelative:     castMoney(v.GetYieldRelative()),
			AccruedInt:        MoneyValue(v.GetAccruedInt()),
			Quantity:          v.GetQuantity(),
			QuantityRest:      v.GetQuantityRest(),
			QuantityDone:      v.GetQuantityDone(),
			CancelDateTime:    v.GetCancelDateTime(),
			CancelReason:      v.GetCancelReason(),
			TradesInfo:        v.GetTradesInfo(),
			AssetUid:          v.GetAssetUid(),
			ChildOperations:   v.GetChildOperations(),
		}
		account.Operations = append(account.Operations, transOperationRet)
	}
}
