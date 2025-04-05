package main

import (
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Account struct {
	Id string
	// Type        pb.AccountType
	Name string
	// Status      int64
	OpenedDate *timestamppb.Timestamp
	ClosedDate *timestamppb.Timestamp
	// AccessLevel pb.AccessLevel
	Portfolio  []PortfolioPosition
	Operations []Operation
}

// Пооучаем список аккаунтов(счетов)
func GetAcc(logger *zap.SugaredLogger, usersService *investgo.UsersServiceClient) map[string]Account {
	accounts := make(map[string]Account)
	status := pb.AccountStatus_ACCOUNT_STATUS_OPEN // ПОтом надо обработать закрытые счета(например ИИС)
	accsResp, err := usersService.GetAccounts(&status)
	if err != nil {
		logger.Errorf(err.Error())
	} else {
		accs := accsResp.GetAccounts()
		for _, acc := range accs {
			account := Account{Id: acc.GetId(),
				// Type:        acc.GetType(),
				Name: acc.GetName(),
				// Status:      int64(acc.GetStatus()),
				OpenedDate: acc.GetOpenedDate(),
				ClosedDate: acc.GetClosedDate(),
				// AccessLevel: acc.GetAccessLevel(),
			}
			accounts[acc.GetId()] = account
			// fmt.Printf("account id = %v\n", accId)
		}
	}

	return accounts

}
