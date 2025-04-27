package tinkoffApi

import (
	"errors"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

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
