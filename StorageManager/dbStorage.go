package StorageManager

import (
	"database/sql"
	"errors"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	_ "golang.org/x/net/context"
)

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

func (ds *DbStorage) Close() error {
	return ds.dbMap.Db.Close()
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
