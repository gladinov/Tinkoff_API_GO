package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Создание БД с именем nameDB
func BuildDB(nameDB string) {
	db, err := sql.Open("sqlite3", nameDB)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Printf("✓ created DB %s\n", nameDB)
}

// Добавление позиций портфеля в БД
func AddPositions(nameDB string, account Account) {
	// Открываем БД
	db, err := sql.Open("sqlite3", nameDB)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Создаем таблицу в БД с динамическим названием портфель_аккаунт
	{
		_, err := db.Exec(PortfolioSqlQuery(account.Id))
		if err != nil {
			panic(err)
		}
		fmt.Printf("✓ created table portfolio_%s\n", account.Id)

	}
	count := 0
	// Добавляем позиции в таблицу БД с динамическим названием портфель_аккаунт
	for _, vals := range account.Portfolio.PortfolioPositios {

		_, err := db.Exec(InsertPortfolioSQL(account.Id),
			vals.Figi,
			vals.InstrumentType,
			vals.Currency,
			vals.Quantity,
			vals.AveragePositionPrice,
			vals.ExpectedYield,
			vals.CurrentNkd,
			vals.CurrentPrice,
			vals.AveragePositionPriceFifo,
			vals.Blocked,
			vals.BlockedLots,
			vals.PositionUid,
			vals.InstrumentUid,
			vals.AssetUid,
			vals.VarMargin,
			vals.ExpectedYieldFifo,
			vals.DailyYield)

		if err != nil {
			panic(err)
		} else {
			count += 1
		}
		// positionId, err := res.LastInsertId()
		// fmt.Printf("added new position: id=%d, error=%v\n", positionId, err)
	}
	fmt.Printf("✓ added %v positions in SQL table portfolio_%s\n", count, account.Id)

}

// Дописать
func AddOperations(nameDB string, account Account) {
	// Открываем БД
	db, err := sql.Open("sqlite3", nameDB)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Создаем таблицу в БД с динамическим названием портфель_аккаунт
	{
		_, err := db.Exec(OperationSqlQuery(account.Id))
		if err != nil {
			panic(err)
		}
		fmt.Printf("✓ created table operations_%s\n", account.Id)
	}
	count := 0
	// Добавляем позиции в таблицу БД с динамическим названием портфель_аккаунт
	for _, vals := range account.Operations {
		_, err := db.Exec(InsertOperationsSQL(account.Id),
			vals.Currency,
			vals.Cursor,
			vals.BrokerAccountId,
			vals.Operation_Id,
			vals.ParentOperationId,
			vals.Name,
			vals.Date.AsTime().Format(time.RFC3339),
			vals.Type,
			vals.Description,
			vals.State,
			vals.InstrumentUid,
			vals.Figi,
			vals.InstrumentType,
			vals.InstrumentKind,
			vals.PositionUid,
			vals.Payment,
			vals.Price,
			vals.Commission,
			vals.Yield,
			vals.YieldRelative,
			vals.AccruedInt,
			vals.Quantity,
			vals.QuantityRest,
			vals.QuantityDone,
			vals.CancelDateTime.AsTime().Format(time.RFC3339),
			vals.CancelReason,
			vals.AssetUid,
		)

		if err != nil {
			panic(err)
		} else {
			count += 1
		}
		// positionId, err := res.LastInsertId()
		// fmt.Printf("added new operations: id=%d, error=%v\n", positionId, err)
	}
	fmt.Printf("✓ added %v operations in SQL table opertions_%s\n", count, account.Id)
}

func AddTickerUid(nameDB string, uidTicker map[string][]string) {
	// Открываем БД
	db, err := sql.Open("sqlite3", nameDB)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Создаем таблицу в БД с динамическим названием портфель_аккаунт
	{
		_, err := db.Exec(`
        DROP TABLE IF EXISTS uidTicker;
        CREATE TABLE IF NOT EXISTS uidTicker (
            id integer primary key,
            instrumentUid                 TEXT,
            ticker                      TEXT,
			class_code	TEXT
        );`)
		if err != nil {
			panic(err)
		}
		fmt.Print("✓ created table uidTicker\n")
	}

	count := 0

	query := `
    insert into uidTicker(instrumentUid, ticker, class_code)
    values (?, ?, ?)
	`

	for key, val := range uidTicker {
		_, err := db.Exec(query, key, val[0], val[1])

		if err != nil {
			panic(err)
		} else {
			count += 1
		}
	}
	fmt.Printf("✓ added %v operations in SQL table uidTicker \n", count)

}

// func AddBondActions(nameDB string, b Bond) {
// 	db, err := sql.Open("sqlite3", nameDB)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer db.Close()

// 	{
// 		_, err := db.Exec()
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Print("✓ created table Amortizations")
// 	}
// }
