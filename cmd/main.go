package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gothanks/myapp/api/tinkoffApi"
	"github.com/gothanks/myapp/internal/database"
	"github.com/gothanks/myapp/internal/service"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

func main() {
	// загружаем конфигурацию для сдк из .yaml файла
	config, err := investgo.LoadConfig("../configs/config.yaml")
	if err != nil {
		log.Fatalf("config loading error %v", err.Error())
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()
	// сдк использует для внутреннего логирования investgo.Logger
	// для примера передадим uber.zap
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"
	l, err := zapConfig.Build()
	logger := l.Sugar()
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Print(err.Error())
		}
	}()
	if err != nil {
		log.Fatalf("logger creating error %v", err)
	}
	// создаем клиента для investAPI, он позволяет создавать нужные сервисы и уже
	// через них вызывать нужные методы
	client, err := investgo.NewClient(ctx, config, logger)
	if err != nil {
		logger.Fatalf("client creating error %v", err.Error())
	}
	defer func() {
		logger.Infof("closing client connection")
		err := client.Stop()
		if err != nil {
			logger.Errorf("client shutdown error %v", err.Error())
		}
	}()

	// Получаем все аккаунты, доступные по данном токену
	// Результатом будет MAP по ключу номер_аккаунта и значениям типа Аккаунт
	accsList, err := tinkoffApi.GetAcc(logger, client)
	if err != nil {
		logger.Errorf("Main. GetAcc error %v", err.Error())
	}
	// opereationsService := client.NewOperationsServiceClient()
	// Создаем БД
	nameDB := "T_API.db"
	database.BuildDB(nameDB)
	// Получаем связку InstrumentUid - AssetUid по всем бумагам в Т-Апи,
	//  т.к. по другому узнать assetUid по конкретной бумаге сервис возможности не дает
	assetUidInstrumentUidMap, err := tinkoffApi.GetAllAssetUids(client)
	if err != nil {
		logger.Errorf("assetUidInstrumentUidMap error %v", err.Error())
	}

	for _, account := range accsList {
		bondReport, err := getBondReports(account, client, assetUidInstrumentUidMap, nameDB)
		if err != nil {
			logger.Errorf("main error %v", err.Error())
		}
		fmt.Println(account.Id)
		fmt.Println()
		fmt.Println(bondReport.BondsInRUB)
		fmt.Println(bondReport.BondsInCNY)
		fmt.Println()
		err = database.AddBondReportsInDB(nameDB, account.Id, bondReport.BondsInRUB)
	}
}

func getBondReports(account tinkoffApi.Account, client *investgo.Client, assetUidInstrumentUidMap map[string]string, nameDB string) (service.Report, error) {
	var bondReport service.Report
	opereationsService := client.NewOperationsServiceClient()
	// Получаем данные по портфелям по кажому счету
	err := tinkoffApi.GetPortf(client, &account)
	if err != nil {
		return bondReport, errors.New("getBondReports" + err.Error())
	}
	// Трансформируем данные портфеля в структуру
	portfolio, err := service.TransPositions(client, &account, assetUidInstrumentUidMap)
	if err != nil {
		return bondReport, errors.New("getBondReports" + err.Error())
	}
	// Добавляем в базу данных
	database.AddPositions(nameDB, account.Id, portfolio.PortfolioPositions)

	// получаем данные по операциям
	err = tinkoffApi.GetOpp(opereationsService, &account)
	if err != nil {
		return bondReport, errors.New("getBondReports" + err.Error())
	}
	// Приводим операции к удобной структуре
	operations := service.TransOperations(account.Operations)
	// добавляем операции в DB
	database.AddOperations(nameDB, account.Id, operations)
	for _, v := range portfolio.BondPositions {
		operationsDb, err := database.GetOperationsFromDBByAssetUid(nameDB, v.Identifiers.AssetUid, account.Id)
		if err != nil {
			return bondReport, errors.New("getBondReports" + err.Error())
		}
		resultBondPosition, err := service.ProcessOperations(client, operationsDb)
		if err != nil {
			return bondReport, errors.New("getBondReports" + err.Error())
		}
		err = bondReport.CreateBondReport(*resultBondPosition)
		if err != nil {
			return bondReport, errors.New("getBondReports" + err.Error())
		}

	}
	return bondReport, nil
}
