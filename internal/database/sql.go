package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gothanks/myapp/internal/service"
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
func AddPositions(nameDB string, accountId string, positions []service.PortfolioPosition) {
	dbPath := fmt.Sprintf("../internal/database/%s", nameDB)
	// Открываем БД
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Создаем таблицу в БД с динамическим названием портфель_аккаунт
	{
		_, err := db.Exec(PortfolioSqlQuery(accountId))
		if err != nil {
			panic(err)
		}
		fmt.Printf("✓ created table portfolio_%s\n", accountId)

	}
	count := 0
	// Добавляем позиции в таблицу БД с динамическим названием портфель_аккаунт
	for _, vals := range positions {

		_, err := db.Exec(InsertPortfolioSQL(accountId),
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
	fmt.Printf("✓ added %v positions in SQL table portfolio_%s\n", count, accountId)

}

// Дописать
func AddOperations(nameDB string, accountId string, operations []service.Operation) {
	dbPath := fmt.Sprintf("../internal/database/%s", nameDB)
	// Открываем БД
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Создаем таблицу в БД с динамическим названием портфель_аккаунт
	{
		_, err := db.Exec(OperationSqlQuery(accountId))
		if err != nil {
			panic(err)
		}
		fmt.Printf("✓ created table operations_%s\n", accountId)
	}
	count := 0
	// Добавляем позиции в таблицу БД с динамическим названием портфель_аккаунт
	for _, vals := range operations {
		_, err := db.Exec(InsertOperationsSQL(accountId),
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
	fmt.Printf("✓ added %v operations in SQL table opertions_%s\n", count, accountId)
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

func GetOperationsFromDBByAssetUid(nameDB, assetUid, accountId string) ([]service.OperationDB, error) {
	dbPath := fmt.Sprintf("../internal/database/%s", nameDB)
	// Открываем БД
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.New("GetOperationsFromDB: sql.Open" + err.Error())
	}
	query := fmt.Sprintf("select name,date, figi, operation_id,quantity_done,instrument_type,instrument_uid,price,currency,accrued_int,commission, payment from operations_%s where asset_uid == '%s'", accountId, assetUid)
	rows, err := db.Query(query)
	if err != nil {
		return nil, errors.New("GetOperationsFromDB: db.Query" + err.Error())
	}
	defer rows.Close()

	operationRes := make([]service.OperationDB, 0)

	for rows.Next() {
		var operation service.OperationDB
		err := rows.Scan(&operation.Name,
			&operation.Date,
			&operation.Figi,
			&operation.Operation_Id,
			&operation.QuantityDone,
			&operation.InstrumentType,
			&operation.InstrumentUid,
			&operation.Price,
			&operation.Currency,
			&operation.AccruedInt,
			&operation.Commission,
			&operation.Payment)
		if err != nil {
			return nil, errors.New("GetOperationsFromDB: rows.Scan" + err.Error())
		}
		operationRes = append(operationRes, operation)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("GetOperationsFromDB: rows.Err()" + err.Error())
	}

	return operationRes, nil
}
