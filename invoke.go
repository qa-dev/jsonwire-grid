package main

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/qa-dev/jsonwire-grid/config"
	"github.com/qa-dev/jsonwire-grid/pool"

	"errors"
	mysqlMigrations "github.com/qa-dev/jsonwire-grid/storage/migrations/mysql"
	"github.com/qa-dev/jsonwire-grid/storage/mysql"
	"github.com/rubenv/sql-migrate"
)

func invokeStorage(config config.Config) (storage pool.StorageInterface, err error) {
	switch config.DB.Implementation {
	case "mysql":
		var db *sqlx.DB
		db, err = sqlx.Open("mysql", config.DB.Connection)
		if err != nil {
			err = errors.New("Database connection error: " + err.Error())
			return
		}

		storage = mysql.NewMysqlStorage(db)

		migrations := &migrate.AssetMigrationSource{
			Asset:    mysqlMigrations.Asset,
			AssetDir: mysqlMigrations.AssetDir,
			Dir:      "storage/migrations/mysql",
		}
		var n int
		n, err = migrate.Exec(db.DB, "mysql", migrations, migrate.Up)
		if err != nil {
			err = errors.New("Migrations failed, " + err.Error())
			return
		}
		log.Infof("Applied %d migrations!\n", n)
	default:
		err = errors.New("Invalid config, unknown param [DB.Implementation=" + config.DB.Implementation + "]")
	}
	return
}
