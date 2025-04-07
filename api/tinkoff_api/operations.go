package tinkoff_api

import (
	"errors"
	"fmt"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

func GetOpp(
	operationsService *investgo.OperationsServiceClient,
	account *Account) error {

	operationsResp, err := operationsService.GetOperationsByCursor(&investgo.GetOperationsByCursorRequest{
		AccountId: account.Id,
		Limit:     1000,
	})
	if err != nil {
		return errors.New("GetOpp: operationsService.GetOperationsByCursor" + err.Error())
	} else {
		ops := operationsResp.GetOperationsByCursorResponse.GetItems()
		account.Operations = append(account.Operations, ops...)
		nextCursor := operationsResp.NextCursor
		for nextCursor != "" {
			operationsResp, err := operationsService.GetOperationsByCursor(&investgo.GetOperationsByCursorRequest{
				AccountId: account.Id,
				Limit:     1000,
				Cursor:    nextCursor,
			})
			if err != nil {
				return errors.New("GetOpp: operationsService.GetOperationsByCursor" + err.Error())
			} else {
				nextCursor = operationsResp.NextCursor
				ops := operationsResp.GetOperationsByCursorResponse.Items
				account.Operations = append(account.Operations, ops...)
			}
		}

	}
	// fmt.Println(account.Operations)
	fmt.Printf("✓ Добавлено %v операции в Account.Operation по счету %s\n", len(account.Operations), account.Id)
	return nil
}
