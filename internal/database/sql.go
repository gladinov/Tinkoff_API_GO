package database

import (
	"database/sql"
	"errors"
	"fmt"

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
	for _, val := range operations {
		_, err := db.Exec(InsertOperationsSQL(accountId),
			val.Currency,
			val.BrokerAccountId,
			val.Operation_Id,
			val.ParentOperationId,
			val.Name,
			val.Date,
			val.Type,
			val.Description,
			val.InstrumentUid,
			val.Figi,
			val.InstrumentType,
			val.InstrumentKind,
			val.PositionUid,
			val.Payment,
			val.Price,
			val.Commission,
			val.Yield,
			val.YieldRelative,
			val.AccruedInt,
			val.QuantityDone,
			val.AssetUid,
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

func GetOperationsFromDBByAssetUid(nameDB, assetUid, accountId string) ([]service.Operation, error) {
	dbPath := fmt.Sprintf("../internal/database/%s", nameDB)
	// Открываем БД
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.New("GetOperationsFromDB: sql.Open" + err.Error())
	}
	query := fmt.Sprintf("select name,date,type, figi, operation_id,quantity_done,instrument_type,instrument_uid,price,currency,accrued_int,commission, payment from operations_%s where asset_uid == '%s' order by date", accountId, assetUid)
	rows, err := db.Query(query)
	if err != nil {
		return nil, errors.New("GetOperationsFromDB: db.Query" + err.Error())
	}
	defer rows.Close()

	operationRes := make([]service.Operation, 0)

	for rows.Next() {
		var operation service.Operation
		err := rows.Scan(&operation.Name,
			&operation.Date,
			&operation.Type,
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

func AddBondReportsInDB(nameDB, accountId string, bondReport []service.BondReport) error {
	if len(bondReport) == 0 {
		return errors.New("Срез пустой")
	}
	dbPath := fmt.Sprintf("../internal/database/%s", nameDB)
	// Открываем БД
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.New("database: AddBondReportsInDB" + err.Error())
	}

	defer db.Close()

	// Создаем таблицу в БД с динамическим названием портфель_аккаунт
	{
		_, err := db.Exec(CreateBondReportQuery(accountId))
		if err != nil {
			return errors.New("database: AddBondReportsInDB" + err.Error())
		}
		fmt.Printf("✓ created table bondReport_%s\n", accountId)
	}
	count := 0
	// Добавляем позиции в таблицу БД с динамическим названием портфель_аккаунт
	for _, report := range bondReport {
		query, err := InsertBondReportsTableQuery(accountId)
		if err != nil {
			return errors.New("database: AddBondReportsInDB" + err.Error())
		}
		_, err = db.Exec(query,
			report.Name,
			report.Ticker,
			report.MaturityDate,
			report.OfferDate,
			report.Duration,
			report.BuyDate,
			report.BuyPrice,
			report.YieldToMaturityOnPurchase,
			report.YieldToOfferOnPurchase,
			report.YieldToMaturity,
			report.YieldToOffer,
			report.CurrentPrice,
			report.Nominal,
			report.Profit,
			report.AnnualizedReturn,
		)

		if err != nil {
			return errors.New("database: AddBondReportsInDB" + err.Error())
		} else {
			count += 1
		}

	}
	fmt.Printf("✓ added %v positions in SQL table bondReport_%s\n", count, accountId)
	return nil
}
