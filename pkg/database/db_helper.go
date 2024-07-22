package database

import (
	"fmt"
	"github.com/go-errors/errors"
	mysqlDriver "gorm.io/driver/mysql"
	postgresDriver "gorm.io/driver/postgres"
	sqlserverDriver "gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// GetGormConnection Возвращает клиент для работы с БД
func GetGormConnection(dbConfig DbConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch dbConfig.Driver {
	case Postgres:
		dialector = postgresDriver.Open(getPostgresConnectionString(dbConfig))
	case MsSql:
		dialector = sqlserverDriver.Open(getMsSqlConnectionString(dbConfig))
	case MySql:
		dialector = mysqlDriver.Open(getMySqlConnectionString(dbConfig))
	}

	if dialector == nil {
		return nil, errors.New("Unknown db driver: " + dbConfig.Driver)
	}

	config := &gorm.Config{}

	if dbConfig.Logging {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	//if dbConfig.DisableAutomaticPing {
	config.DisableAutomaticPing = true
	//}

	gormDb, err := gorm.Open(dialector, config)

	if err != nil {
		return nil, err
	}

	db, err := gormDb.DB()

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Hour)

	if dbConfig.MaxOpenConnections > 0 {
		db.SetMaxOpenConns(dbConfig.MaxOpenConnections)
	}

	if dbConfig.MaxIdleConnections > 0 {
		db.SetMaxIdleConns(dbConfig.MaxIdleConnections)
	}

	return gormDb, nil
}

// getPostgresConnectionString Возвращает строку (DSN) для создания соединения с Postgres
func getPostgresConnectionString(dbConfig DbConfig) string {
	sslMode := "disable"

	if dbConfig.SslMode {
		sslMode = "enable"
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Db, sslMode)
}

// getMsSqlConnectionString Возвращает строку (DSN) для создания соединения с MSSQL
func getMsSqlConnectionString(dbConfig DbConfig) string {
	return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Port, dbConfig.Db)
}

// getMySqlConnectionString Возвращает строку (DSN) для создания соединения с Mysql
func getMySqlConnectionString(dbConfig DbConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Db)
}
