package mysql

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/qa-dev/jsonwire-grid/config"
	"github.com/qa-dev/jsonwire-grid/pool"
	mysqlMigrations "github.com/qa-dev/jsonwire-grid/storage/migrations/mysql"
	"github.com/rubenv/sql-migrate"
)

type Factory struct {
}

func (f *Factory) Create(config config.Config) (pool.StorageInterface, error) {
	db, err := sqlx.Open("mysql", config.DB.Connection)
	if err != nil {
		err = errors.New("Database connection error: " + err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		err = errors.New("Database connection not establish: " + err.Error())
		return nil, err
	}

	db.SetMaxIdleConns(0)  // this is the root problem! set it to 0 to remove all idle connections
	db.SetMaxOpenConns(10) // or whatever is appropriate for your setup.

	storage := NewMysqlStorage(db)

	migrations := &migrate.AssetMigrationSource{
		Asset:    mysqlMigrations.Asset,
		AssetDir: mysqlMigrations.AssetDir,
		Dir:      "storage/migrations/mysql",
	}
	n, err := migrate.Exec(db.DB, "mysql", migrations, migrate.Up)
	if err != nil {
		err = errors.New("Migrations failed, " + err.Error())
		return nil, err
	}
	log.Infof("Applied %d migrations!\n", n)
	return storage, nil
}
