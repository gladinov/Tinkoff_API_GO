package tinkoff_api

import (
	"errors"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func GetPortf(client *investgo.Client,
	account *Account) error {
	operationsService := client.NewOperationsServiceClient()
	id := account.Id
	portfolioResp, err := operationsService.GetPortfolio(id,
		pb.PortfolioRequest_RUB)
	if err != nil {
		return errors.New("GetPortf: operationsService.GetPortfolio" + err.Error())
	}
	account.Portfolio = portfolioResp.GetPositions()

	return nil
}
