package main

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
)

/*type DBConnection interface {
		Username string
		Password string
		Hostname string
		Port     string
		DBName   string
		Dialect  string

}*/

type Storage struct {
	DBMap *gorp.DbMap
}

var storageInstance *Storage

func GetStorageInstance(dialect, hostname, port, DBName, username, password string) *Storage {
	if storageInstance == nil {
		storageInstance = NewStorage(dialect, hostname, port, DBName, username, password)
	}

	return storageInstance
}

func NewStorage(dialect, hostname, port, DBName, username, password string) *Storage {
	storageInstance := Storage{}
	var connectionString = fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s", username, DBName, password, hostname, port)

	db, err := sql.Open(dialect, connectionString)
	if err != nil {
		log.Fatalf("[DEBUG]: DB error DB driver: %+v", err)
	}

	storageInstance.DBMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	storageInstance.DBMap.AddTableWithName(Price{}, "items").SetKeys(true, "Id")

	err = storageInstance.DBMap.CreateTablesIfNotExists()
	checkErr(err, "[DEBUG]: DB Create tables failed")

	// Delete any existing rows
	err = storageInstance.DBMap.TruncateTables()
	checkErr(err, "[DEBUG]: DB TruncateTables error")

	return &storageInstance
}

func (storage *Storage) Write(priceList PriceList) (n int, err error) {

	//defer timeTrack(time.Now(), "[TIMER] insert to DB")

	trans, err := storage.DBMap.Begin()
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

func (storage *Storage) Close() error {
	return storage.DBMap.Db.Close()
}
