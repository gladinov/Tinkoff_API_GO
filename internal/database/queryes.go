package database

import (
	"errors"
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
        broker_account_id       TEXT,
        operation_id            TEXT,
        parent_operation_id     TEXT,
        name                    TEXT,
        date                    DATETIME, 
        type                    INTEGER,
        description             TEXT,
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
        broker_account_id,
        operation_id,
        parent_operation_id,
        name,
        date,
        type,
        description,
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
        asset_uid
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `, accountId)

	return InsertOperationQuery
}

func CreateBondReportQuery(accountId string) (string, error) {
	if !isValidAccountId(accountId) {
		return "", errors.New("database: CreateBondReportQuery: Недопустимое имя аккаунта")
	}

	createBondReportsTableQuery := fmt.Sprintf(`
    DROP TABLE IF EXISTS bond_reports_%s;
    CREATE TABLE IF NOT EXISTS bond_reports_%s (
        id integer primary key,
        name                            TEXT,
        ticker                          TEXT,
        maturity_date                   DATETIME,
        offer_date                      DATETIME,
        duration                        INTEGER,
        buy_date                        DATETIME,
        buy_price                       REAL,
        yield_to_maturity_on_purchase   REAL,
        yield_to_offer_on_purchase      REAL,
        yield_to_maturity               REAL,
        yield_to_offer                  REAL,
        current_price                   REAL,
        nominal                         REAL,
        profit                          REAL,
        annualized_return               REAL
    );`, accountId, accountId)

	return createBondReportsTableQuery, nil
}

func InsertBondReportsTableQuery(accountId string) (string, error) {
	if !isValidAccountId(accountId) {
		return "", errors.New("database: CreateBondReportQuery: Недопустимое имя аккаунта")
	}

	InsertBondReportQuery := fmt.Sprintf(`
    INSERT INTO bond_reports_%s (
        name,
        ticker,
        maturity_date,
        offer_date,
        duration,
        buy_date,
        buy_price,
        yield_to_maturity_on_purchase,
        yield_to_offer_on_purchase,
        yield_to_maturity,
        yield_to_offer,
        current_price,
        nominal,
        profit,
        annualized_return
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`, accountId)

	return InsertBondReportQuery, nil
}
