package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
)

type Storage interface {
	Write([]*interface{}) (int, error)
	Truncate() error
	Drop() error
	Close() error
}

type DbStorage struct {
	dbMap *gorp.DbMap
}

var dbStorageInst *DbStorage

func GetDbStorageInst(dialectName, dsn, tableName string, tableStruct interface{}) (*DbStorage, error) {
	if dbStorageInst == nil {
		return NewDbStorage(dialectName, dsn, tableName, tableStruct)
	}
	return dbStorageInst, nil
}

func NewDbStorage(dialectName, dsn, tableName string, tableStruct interface{}) (*DbStorage, error) {

	db, err := sql.Open(dialectName, dsn)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	dbStorageInst := DbStorage{}
	dialect, err := dbStorageInst.getDialectByName(dialectName)
	if err != nil {
		return nil, err
	}

	dbStorageInst.dbMap = &gorp.DbMap{Db: db, Dialect: dialect}

	dbStorageInst.dbMap.AddTableWithName(tableStruct, tableName).SetKeys(true, "Id")

	err = dbStorageInst.dbMap.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}

	return &dbStorageInst, nil
}

func (ds DbStorage) getDialectByName(dialectName string) (gorp.Dialect, error) {
	switch dialectName {
	case "postgres":
		return gorp.PostgresDialect{}, nil
	case "mysql":
		return gorp.MySQLDialect{}, nil
	case "sqlite3":
		return gorp.SqliteDialect{}, nil
	case "oracle":
		return gorp.OracleDialect{}, nil
	case "sqlserver":
		return gorp.SqlServerDialect{}, nil
	}

	return nil, errors.New("No dialect associated with provided name")
}

func (ds *DbStorage) Close() (err error) {
	err = ds.dbMap.Db.Close()
	return err
}

func (ds *DbStorage) Truncate() (err error) {
	err = ds.dbMap.TruncateTables()
	return err
}

func (ds *DbStorage) Drop() (err error) {
	err = ds.dbMap.DropTablesIfExists()
	return err
}

func (ds *DbStorage) Write(dataSlice []*interface{}) (int, error) {
	// Begin transaction
	trans, err := ds.dbMap.Begin()
	if err != nil {
		return 0, err
	}

	for _, item := range dataSlice {
		trans.Insert(item)
	}

	if err = trans.Commit(); err != nil {
		return 0, err
	}

	return len(dataSlice), nil
}

type Stor struct {
	DBMap *gorp.DbMap
}

var storInst *Stor

func GetStorInstance(dbCon StorageConfig, tableName string, tableStruct interface{}) *Stor {
	if storInst == nil {
		storInst = NewStorage(dbCon, tableName, tableStruct)
	}

	return storInst
}

func NewStorage(dbCon StorageConfig, tableName string, tableStruct interface{}) *Stor {
	storageInst := Stor{}
	dbConnectionString := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s", dbCon.Username, dbCon.DBName, dbCon.Password, dbCon.Hostname, dbCon.Port)

	db, err := sql.Open(dbCon.Type, dbConnectionString)
	if err != nil {
		log.Fatalf("[DEBUG]: DB error DB driver: %+v", err)
	}

	storageInst.DBMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	storageInst.DBMap.AddTableWithName(tableStruct, tableName).SetKeys(true, "Id")

	err = storageInst.DBMap.CreateTablesIfNotExists()
	checkErr(err, "[DEBUG]: DB Create tables failed")

	// Delete any existing rows
	err = storageInst.DBMap.TruncateTables()
	checkErr(err, "[DEBUG]: DB TruncateTables error")

	return &storageInst
}

func (st *Stor) Write(priceList PriceCollection) (n int, err error) {

	trans, err := st.DBMap.Begin()
	if err != nil {
		return 0, err
	}

	for _, price := range priceList {
		trans.Insert(price)
	}

	if err = trans.Commit(); err != nil {
		return 0, err
	}

	return len(priceList), nil
}

func (st *Stor) Close() error {
	return st.DBMap.Db.Close()
}
