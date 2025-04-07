package service

import (
	"errors"
	"fmt"

	"github.com/gothanks/myapp/api/tinkoff_api"
	"github.com/gothanks/myapp/other_func"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

type Portfolio struct {
	PortfolioPositions []PortfolioPosition
	BondPosions        []Bond
}

type PortfolioPosition struct {
	AccountId                string
	Figi                     string
	InstrumentType           string
	Currency                 string
	Quantity                 float64
	AveragePositionPrice     float64
	ExpectedYield            float64
	CurrentNkd               float64
	CurrentPrice             float64
	AveragePositionPriceFifo float64
	Blocked                  bool
	BlockedLots              float64
	PositionUid              string
	InstrumentUid            string
	AssetUid                 string
	VarMargin                float64
	ExpectedYieldFifo        float64
	DailyYield               float64
}

// Обрабатываем в нормальный формат портфеля
func TransPositions(client *investgo.Client,
	account *tinkoff_api.Account, assetUidInstrumentUidMap map[string]string) Portfolio {
	Portfolio := Portfolio{}
	for _, v := range account.Portfolio {
		if v.InstrumentType == "bond" {
			BondPosition := Bond{
				Identifiers: Identifiers{
					Figi:          v.GetFigi(),
					InstrumentUid: v.GetInstrumentUid(),
					PositionUid:   v.GetPositionUid(),
				},
				InstrumentType:           v.GetInstrumentType(),
				Currency:                 v.GetAveragePositionPrice().Currency,
				Quantity:                 other_func.CastMoney(v.GetQuantity()),
				AveragePositionPrice:     other_func.MoneyValue(v.GetAveragePositionPrice()),
				ExpectedYield:            other_func.CastMoney(v.GetExpectedYield()),
				CurrentNkd:               other_func.MoneyValue(v.GetCurrentNkd()),
				CurrentPrice:             other_func.MoneyValue(v.GetCurrentPrice()),
				AveragePositionPriceFifo: other_func.MoneyValue(v.GetAveragePositionPriceFifo()),
				Blocked:                  v.GetBlocked(),
				ExpectedYieldFifo:        other_func.CastMoney(v.GetExpectedYieldFifo()),
				DailyYield:               other_func.MoneyValue(v.GetDailyYield()),
			}
			// Получаем AssetUid с помощью МАПЫ assetUidInstrumentUidMap
			BondPosition.Identifiers.AssetUid = assetUidInstrumentUidMap[BondPosition.Identifiers.InstrumentUid]
			// Получаем Тикер, Режим торгов и Короткое имя инструмента
			BondPosition.GetBondsActionsFromPortfolio(client)
			//  Получение данных с московской биржи
			BondPosition.GetActionFromMoex()
			Portfolio.BondPosions = append(Portfolio.BondPosions, BondPosition)
		} else {
			transPosionRet := PortfolioPosition{
				Figi:                     v.GetFigi(),
				InstrumentType:           v.GetInstrumentType(),
				Currency:                 v.GetAveragePositionPrice().Currency,
				Quantity:                 other_func.CastMoney(v.GetQuantity()),
				AveragePositionPrice:     other_func.MoneyValue(v.GetAveragePositionPrice()),
				ExpectedYield:            other_func.CastMoney(v.GetExpectedYield()),
				CurrentNkd:               other_func.MoneyValue(v.GetCurrentNkd()),
				CurrentPrice:             other_func.MoneyValue(v.GetCurrentPrice()),
				AveragePositionPriceFifo: other_func.MoneyValue(v.GetAveragePositionPriceFifo()),
				Blocked:                  v.GetBlocked(),
				BlockedLots:              other_func.CastMoney(v.GetBlockedLots()),
				PositionUid:              v.GetPositionUid(),
				InstrumentUid:            v.GetInstrumentUid(),
				AssetUid:                 assetUidInstrumentUidMap[v.GetInstrumentUid()],
				VarMargin:                other_func.MoneyValue(v.GetVarMargin()),
				ExpectedYieldFifo:        other_func.CastMoney(v.GetExpectedYieldFifo()),
				DailyYield:               other_func.MoneyValue(v.GetDailyYield()),
			}
			Portfolio.PortfolioPositions = append(Portfolio.PortfolioPositions, transPosionRet)
		}
	}
	fmt.Printf("✓ Добавлено %v позиций в Account.PortfolioPositions по счету %s\n", len(Portfolio.PortfolioPositions), account.Id)
	fmt.Printf("✓ Добавлено %v позиций в Account.PortfolioBondPositions по счету %s\n", len(Portfolio.BondPosions), account.Id)
	return Portfolio
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
