package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gothanks/myapp/api/tinkoff_api"
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
			log.Printf(err.Error())
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
	accsList, err := tinkoff_api.GetAcc(logger, client)
	if err != nil {
		logger.Errorf("Main. GetAcc error %v", err.Error())
	}
	opereationsService := client.NewOperationsServiceClient()
	// Создаем БД
	nameDB := "T_API.db"
	database.BuildDB(nameDB)
	// Получаем связку InstrumentUid - AssetUid по всем бумагам в Т-Апи,
	//  т.к. по другому узнать assetUid по конкретной бумаге сервис возможности не дает
	assetUidInstrumentUidMap, err := service.GetAllAssetUids(client)
	if err != nil {
		logger.Errorf("assetUidInstrumentUidMap error %v", err.Error())
	}

	for _, account := range accsList {
		// Запросы в tinkoff.Api

		// Получаем данные по портфелям по кажому счету
		err := tinkoff_api.GetPortf(client, &account)
		if err != nil {
			logger.Errorf("tinkoff_api.GetPortf error %v", err.Error())
		}
		// Трансформируем данные портфеля в структуру
		portfolio := service.TransPositions(client, &account, assetUidInstrumentUidMap)
		// Добавляем в базу данных
		database.AddPositions(nameDB, account.Id, portfolio.PortfolioPositions)
		// получаем данные по операциям
		err = tinkoff_api.GetOpp(opereationsService, &account)
		if err != nil {
			logger.Errorf("tinkoff_api.GetOpp error %v", err.Error())
		}
		// Приводим операции к удобной структуре
		operations := service.TransOperations(account.Operations)
		// добавляем операции в DB
		database.AddOperations(nameDB, account.Id, operations)
		for _, v := range portfolio.BondPosions {
			fmt.Println()
			fmt.Println(v.Name)
			fmt.Println()
			fmt.Println(database.GetOperationsFromDBByAssetUid(nameDB, v.Identifiers.AssetUid, account.Id))
		}
	}
}
