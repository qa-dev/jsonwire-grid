package tests

import (
	"crypto/rand"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/storage/mysql"
	"github.com/rubenv/sql-migrate"
	"os"
	"strings"
	"testing"
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
	const nameLen = 32
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, nameLen)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = chars[v%byte(len(chars))]
	}
	dbName := "TEST_" + string(bytes)
	_, err := commonDBConnection.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		err = errors.New("Database create error: " + err.Error())
	}
	if err != nil {
		panic(err.Error())
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

	storage := mysql.NewMysqlStorage(db)
	if err != nil {
		panic("Error initialisation storage: " + err.Error())
	}
	deferFunc := func() {
		commonDBConnection.Exec("DROP DATABASE " + dbName)
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
	testStorage_Add(t, mv)
}

// TestMysqlStorage_Add_Repeat see testStorage_Add_Repeat
func TestMysqlStorage_Add_Repeat(t *testing.T) {
	testStorage_Add_Repeat(t, mv)

}

// TestMysqlStorage_GetAll see testStorage_GetAll
func TestMysqlStorage_GetAll(t *testing.T) {
	testStorage_GetAll(t, mv)
}

// TestMysqlStorage_GetByAddress see testStorage_GetByAddress
func TestMysqlStorage_GetByAddress(t *testing.T) {
	testStorage_GetByAddress(t, mv)
}

// TestMysqlStorage_GetBySession see testStorage_GetBySession
func TestMysqlStorage_GetBySession(t *testing.T) {
	testStorage_GetBySession(t, mv)
}

// TestMysqlStorage_GetCountWithStatus see testStorage_GetCountWithStatus
func TestMysqlStorage_GetCountWithStatus(t *testing.T) {
	testStorage_GetCountWithStatus(t, mv)
}

// TestMysqlStorage_Remove see testStorage_Remove
func TestMysqlStorage_Remove(t *testing.T) {
	testStorage_Remove(t, mv)
}

// TestMysqlStorage_ReserveAvailable_Positive_Filtration see testStorage_ReserveAvailable_Positive_Filtration
func TestMysqlStorage_ReserveAvailable_Positive_Filtration(t *testing.T) {
	testStorage_ReserveAvailable_Positive_Filtration(t, mv)
}

// TestMysqlStorage_ReserveAvailable_Positive_NoFiltration see testStorage_ReserveAvailable_Positive_NoFiltration
func TestMysqlStorage_ReserveAvailable_Positive_NoFiltration(t *testing.T) {
	testStorage_ReserveAvailable_Positive_NoFiltration(t, mv)
}

// TestMysqlStorage_ReserveAvailable_Negative see testStorage_ReserveAvailable_Negative
func TestMysqlStorage_ReserveAvailable_Negative(t *testing.T) {
	testStorage_ReserveAvailable_Negative(t, mv)
}

// TestMysqlStorage_SetAvailable see testStorage_SetAvailable
func TestMysqlStorage_SetAvailable(t *testing.T) {
	testStorage_SetAvailable(t, mv)
}

// TestMysqlStorage_SetBusy see testStorage_SetBusy
func TestMysqlStorage_SetBusy(t *testing.T) {
	testStorage_SetBusy(t, mv)
}
