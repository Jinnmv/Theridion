package main

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
)

type Mapper struct {
	DBMap *gorp.DbMap
}

type ModelsCatalog struct {
	Id           uint64 `db:"id"`
	Name         string `db:"name"`
	Category     string `db:"category"`
	SubCategory  string `db:"sub_category"`
	Manufacturer string `db:"manufacturer"`
	Scale        string `db:"scale"`
	Sku          string `db:"sku"`
}

func NewMapper(dbCon StorageConfig) *Mapper {

	mapper := Mapper{}
	dbConnectionString := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s", dbCon.Username, dbCon.DBName, dbCon.Password, dbCon.Hostname, dbCon.Port)

	db, err := sql.Open(dbCon.Type, dbConnectionString)
	if err != nil {
		log.Fatalf("[DEBUG]: DB error DB driver: %+v", err)
	}

	mapper.DBMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	mapper.DBMap.AddTableWithName(Price{}, "catalog").SetKeys(true, "Id")

	err = mapper.DBMap.CreateTablesIfNotExists()
	checkErr(err, "[DEBUG]: DB Create tables failed")

	return &mapper
}
