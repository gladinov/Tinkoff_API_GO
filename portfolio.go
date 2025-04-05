package main

import (
	"errors"
	"fmt"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
)

type Portfolio struct {
	PortfolioPositios []PortfolioPosition
}

type PortfolioPosition struct {
	accoutId                 string
	Figi                     string
	InstrumentType           string
	Currency                 string
	Quantity                 float64
	AveragePositionPrice     float64
	ExpectedYield            float64
	CurrentNkd               float64
	AveragePositionPricePt   float64
	CurrentPrice             float64
	AveragePositionPriceFifo float64
	QuantityLots             float64
	Blocked                  bool
	BlockedLots              float64
	PositionUid              string
	InstrumentUid            string
	AssetUid                 string
	VarMargin                float64
	ExpectedYieldFifo        float64
	DailyYield               float64
}

func GetPortf(logger *zap.SugaredLogger,
	client *investgo.Client,
	account *Account) error {
	operationsService := client.NewOperationsServiceClient()
	id := account.Id
	portfolioResp, err := operationsService.GetPortfolio(id,
		pb.PortfolioRequest_RUB)
	if err != nil {
		return errors.New("GetPortf: operationsService.GetPortfolio" + err.Error())
	}
	positions := portfolioResp.GetPositions()
	assetUidInstrumentUidMap, err := GetAllAssetUids(client)
	if err != nil {
		return errors.New("GetPortf: GetAllAssetUids" + err.Error())
	}
	transPositions(positions, account, assetUidInstrumentUidMap)
	return nil
}

// Обрабатываем в нормальный формат портфеля
func transPositions(positions []*pb.PortfolioPosition, account *Account, assetUidInstrumentUidMap map[string]string) {
	for _, v := range positions {
		// fmt.Println(v.GetCurrentNkd())
		transPosionRet := PortfolioPosition{
			accoutId:                 account.Id,
			Figi:                     v.GetFigi(),
			InstrumentType:           v.GetInstrumentType(),
			Currency:                 v.GetAveragePositionPrice().Currency,
			Quantity:                 castMoney(v.GetQuantity()),
			AveragePositionPrice:     MoneyValue(v.GetAveragePositionPrice()),
			ExpectedYield:            castMoney(v.GetExpectedYield()),
			CurrentNkd:               MoneyValue(v.GetCurrentNkd()),
			CurrentPrice:             MoneyValue(v.GetCurrentPrice()),
			AveragePositionPriceFifo: MoneyValue(v.GetAveragePositionPriceFifo()),
			Blocked:                  v.GetBlocked(),
			BlockedLots:              castMoney(v.GetBlockedLots()),
			PositionUid:              v.GetPositionUid(),
			InstrumentUid:            v.GetInstrumentUid(),
			AssetUid:                 assetUidInstrumentUidMap[v.GetInstrumentUid()],
			VarMargin:                MoneyValue(v.GetVarMargin()),
			ExpectedYieldFifo:        castMoney(v.GetExpectedYieldFifo()),
			DailyYield:               MoneyValue(v.GetDailyYield()),
		}
		account.Portfolio.PortfolioPositios = append(account.Portfolio.PortfolioPositios, transPosionRet)
	}
	fmt.Printf("✓ Добавлено %v позиций в Account.Portfolio по счету %s\n", len(account.Portfolio.PortfolioPositios), account.Id)
}

func GetAllAssetUids(client *investgo.Client) (map[string]string, error) {
	instrumentService := client.NewInstrumentsServiceClient()
	answer, err := instrumentService.GetAssets()
	if err != nil {
		return nil, errors.New("GetAllAssetUids: instrumentService.GetAssets" + err.Error())
	}
	assetUidInstrumentUidMap := make(map[string]string)
	for _, v := range answer.AssetsResponse.Assets {
		asset_uid := v.Uid

		for _, instrument := range v.Instruments {
			instrument_uid := instrument.Uid
			assetUidInstrumentUidMap[instrument_uid] = asset_uid
		}
	}
	return assetUidInstrumentUidMap, nil
}
