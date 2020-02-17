package mysql

import (
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/storage/tests"
)

const (
	dbNamePlaceholder             = "%dbname%"
	dbConnectionStringEnvVariable = "TEST_MYSQL_CONNECTION"
)

var (
	connectionString   string
	commonDBConnection *sqlx.DB
	mv                 PrepareMysql
)

type PrepareMysql struct{}

func (p PrepareMysql) SetUp() {
	connectionString = os.Getenv(dbConnectionStringEnvVariable)
	if !strings.Contains(connectionString, dbNamePlaceholder) {
		panic(fmt.Sprintf("Pleas use placeholder %s in %s", dbNamePlaceholder, dbConnectionStringEnvVariable))
	}

	db, err := sqlx.Open("mysql", strings.Replace(connectionString, dbNamePlaceholder, "", 1))
	if err != nil {
		panic("Database connection error: " + err.Error())
	}
	commonDBConnection = db
}

func (p PrepareMysql) TearDown() {
	commonDBConnection.Close()
}

func (p PrepareMysql) CreateStorage() (pool.StorageInterface, func()) {

	dbName := tests.CreateDbName()
	_, err := commonDBConnection.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		panic("Database create error: " + err.Error())
	}

	replace := strings.Replace(connectionString, dbNamePlaceholder, dbName, 1)
	db, err := sqlx.Open("mysql", replace)
	if err != nil {
		panic("Database connection error: " + err.Error())
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "../migrations/mysql",
	}
	_, err = migrate.Exec(db.DB, "mysql", migrations, migrate.Up)
	if err != nil {
		panic("Migrations failed, " + err.Error())
	}

	storage := NewMysqlStorage(db)
	if err != nil {
		panic("Error initialisation storage: " + err.Error())
	}
	deferFunc := func() {
		_, _ = commonDBConnection.Exec("DROP DATABASE " + dbName)
		db.Close()
	}
	return storage, deferFunc
}

func TestMain(m *testing.M) {
	mv = PrepareMysql{}
	mv.SetUp()
	retCode := m.Run()
	mv.TearDown()
	os.Exit(retCode)
}

// TestMysqlStorage_Add see testStorage_Add
func TestMysqlStorage_Add(t *testing.T) {
	tests.TestStorage_Add(t, mv)
}

// TestMysqlStorage_Add_Repeat see testStorage_Add_Repeat
func TestMysqlStorage_Add_Repeat(t *testing.T) {
	tests.TestStorage_Add_Repeat(t, mv)

}

// TestStorage_Add_Limit_Overflow see testStorage_Add_Limit_Overflow
func TestStorage_Add_Limit_Overflow(t *testing.T) {
	tests.TestStorage_Add_Limit_Overflow(t, mv)

}

// TestMysqlStorage_GetAll see testStorage_GetAll
func TestMysqlStorage_GetAll(t *testing.T) {
	tests.TestStorage_GetAll(t, mv)
}

// TestMysqlStorage_GetByAddress see testStorage_GetByAddress
func TestMysqlStorage_GetByAddress(t *testing.T) {
	tests.TestStorage_GetByAddress(t, mv)
}

// TestMysqlStorage_GetBySession see testStorage_GetBySession
func TestMysqlStorage_GetBySession(t *testing.T) {
	tests.TestStorage_GetBySession(t, mv)
}

// TestMysqlStorage_GetCountWithStatus see testStorage_GetCountWithStatus
func TestMysqlStorage_GetCountWithStatus(t *testing.T) {
	tests.TestStorage_GetCountWithStatus(t, mv)
}

// TestMysqlStorage_Remove see testStorage_Remove
func TestMysqlStorage_Remove(t *testing.T) {
	tests.TestStorage_Remove(t, mv)
}

// TestMysqlStorage_ReserveAvailable_Positive see testStorage_ReserveAvailable_Positive
func TestMysqlStorage_ReserveAvailable_Positive(t *testing.T) {
	tests.TestStorage_ReserveAvailable_Positive(t, mv)
}

// TestMysqlStorage_ReserveAvailable_Negative see testStorage_ReserveAvailable_Negative
func TestMysqlStorage_ReserveAvailable_Negative(t *testing.T) {
	tests.TestStorage_ReserveAvailable_Negative(t, mv)
}

// TestMysqlStorage_SetAvailable see testStorage_SetAvailable
func TestMysqlStorage_SetAvailable(t *testing.T) {
	tests.TestStorage_SetAvailable(t, mv)
}

// TestMysqlStorage_SetBusy see testStorage_SetBusy
func TestMysqlStorage_SetBusy(t *testing.T) {
	tests.TestStorage_SetBusy(t, mv)
}

// TestMysqlStorage_UpdateAdderss_UpdatesValue see testStorage_UpdateAdderss_UpdatesValue
func TestMysqlStorage_UpdateAdderss_UpdatesValue(t *testing.T) {
	tests.TestStorage_UpdateAdderss_UpdatesValue(t, mv)
}

// TestMysqlStorage_UpdateAdderss_ReturnsErrNotFound see testStorage_UpdateAdderss_ReturnsErrNotFound
func TestMysqlStorage_UpdateAdderss_ReturnsErrNotFound(t *testing.T) {
	tests.TestStorage_UpdateAdderss_ReturnsErrNotFound(t, mv)
}
