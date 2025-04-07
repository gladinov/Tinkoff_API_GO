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

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
)

func main() {

	// загружаем конфигурацию для сдк из .yaml файла
	config, err := investgo.LoadConfig("config.yaml")
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

	// создаем клиента для сервиса счетов
	// Получаем все аккаунты, доступные по данном токену
	usersService := client.NewUsersServiceClient()

	// Результатом будет MAP по ключу номер_аккаунта и значениям типа Аккаунт
	accsList := GetAcc(logger, usersService)

	opereationsService := client.NewOperationsServiceClient()
	//Создаем подключение к инструментам API

	// Создаем БД
	nameDB := "T_API.db"
	BuildDB(nameDB)
	// //////

	for _, account := range accsList {
		// Получаем данные по портфелям по кажому счету
		GetPortf(logger, client, &account)
		// Добавляем в базу данных
		// AddPositions(nameDB, account)
		// получаем данные по операциям
		GetOpp(logger, opereationsService, &account)
		// добавляем операции в DB
		// AddOperations(nameDB, account)
		for _, v := range account.Portfolio.BondPosions {
			fmt.Println(v)
			fmt.Println()
		}
	}

	// пробная версия получения тикера по uid

}
