package database

import (
	"fmt"
	"log"
	"regexp"
)

// Функция для проверки имени аккаунта. Нейронки говорят , что помогает от SQL иньекций. Пусть будет, пока не так важно.
func isValidAccountId(accountId string) bool {
	validId := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validId.MatchString(accountId)
}

// Динамический Query для создания таблицы портфель_аккаунт
func PortfolioSqlQuery(accountId string) string {
	if !isValidAccountId(accountId) {
		log.Fatal("Недопустимое имя аккаунта")
	}
	PortfolioQuery := fmt.Sprintf(`
        DROP TABLE IF EXISTS portfolio_%s;
        CREATE TABLE IF NOT EXISTS portfolio_%s (
            id integer primary key,
            figi                      TEXT,
            instrumentType            TEXT,
            currency                  TEXT,
            quantity                  REAL,
            averagePositionPrice      REAL,
            expectedYield             REAL,
            currentNkd                REAL,
            currentPrice              REAL,
            averagePositionPriceFifo  REAL,
            blocked                   BOOLEAN,
            blockedLots               REAL,
            positionUid               TEXT,
            instrumentUid             TEXT,
            asset_uid                 TEXT,
            varMargin                 REAL,
            expectedYieldFifo         REAL,
            dailyYield                REAL
        );`, accountId, accountId)

	return PortfolioQuery
}

// Динамический Query для добавления значений в таблицу портфель_аккаунт
func InsertPortfolioSQL(accountId string) string {
	if !isValidAccountId(accountId) {
		log.Fatal("Недопустимое имя аккаунта")
	}

	InsertPortfolioQuery := fmt.Sprintf(`
    INSERT INTO portfolio_%s (
        figi,
        instrumentType,
        currency,
        quantity,
        averagePositionPrice,
        expectedYield,
        currentNkd,
        currentPrice,
        averagePositionPriceFifo,
        blocked,
        blockedLots,
        positionUid,
        instrumentUid,
        asset_uid,
        varMargin,
        expectedYieldFifo,
        dailyYield
    	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, accountId)

	return InsertPortfolioQuery
}

func OperationSqlQuery(accountId string) string {
	if !isValidAccountId(accountId) {
		log.Fatal("Недопустимое имя аккаунта")
	}
	CreateOperationTableQuery := fmt.Sprintf(`
    DROP TABLE IF EXISTS operations_%s;
    CREATE TABLE IF NOT EXISTS operations_%s (
        id integer primary key,
        currency                TEXT,
        cursor                  TEXT,
        broker_account_id       TEXT,
        operation_id            TEXT,
        parent_operation_id     TEXT,
        name                    TEXT,
        date                    TEXT, 
        type                    INTEGER,
        description             TEXT,
        state                   INTEGER,
        instrument_uid          TEXT,
        figi                    TEXT,
        instrument_type         TEXT,
        instrument_kind         TEXT,
        position_uid            TEXT,
        payment                 REAL,
        price                   REAL,
        commission              REAL,
        yield                   REAL,
        yield_relative          REAL,
        accrued_int             REAL,
        quantity_done           INTEGER,
        cancel_date_time        TEXT, 
        cancel_reason           TEXT,
        asset_uid               TEXT
    );`, accountId, accountId)

	return CreateOperationTableQuery
}

func InsertOperationsSQL(accountId string) string {
	if !isValidAccountId(accountId) {
		log.Fatal("Недопустимое имя аккаунта")
	}

	InsertOperationQuery := fmt.Sprintf(`
    INSERT INTO operations_%s (
        currency,
        cursor,
        broker_account_id,
        operation_id,
        parent_operation_id,
        name,
        date,
        type,
        description,
        state,
        instrument_uid,
        figi,
        instrument_type,
        instrument_kind,
        position_uid,
        payment,
        price,
        commission,
        yield,
        yield_relative,
        accrued_int,
        quantity_done,
        cancel_date_time,
        cancel_reason,
        asset_uid
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `, accountId)

	return InsertOperationQuery
}

// func InsertActionsMoexSqlQuery() string {
// 	PortfolioQuery := `
//         DROP TABLE IF EXISTS amortization;
//         CREATE TABLE IF NOT EXISTS amortization(
//             id integer primary key,
//             isin                 TEXT,
//             name                      TEXT,
//             issuevalue            TEXT,
//             coupondate                  TEXT,
//             recorddate                  REAL,
//             startdate      REAL,
//             initialfacevalue             REAL,
//             facevalue                REAL,
//             faceunit    REAL,
//             currentPrice              REAL,
//             averagePositionPriceFifo  REAL,
//             quantityLots              REAL,
//             blocked                   BOOLEAN,
//             blockedLots               REAL,
//             positionUid               TEXT,
//             instrumentUid             TEXT,
//             varMargin                 REAL,
//             expectedYieldFifo         REAL,
//             dailyYield                REAL
//         );`

// 	return PortfolioQuery
// }
