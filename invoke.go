package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/qa-dev/jsonwire-grid/config"

	"errors"
	"github.com/qa-dev/jsonwire-grid/storage"
	"github.com/qa-dev/jsonwire-grid/storage/mysql"
)

func invokeStorageFactory(config config.Config) (factory storage.StorageFactoryInterface, err error) {
	switch config.DB.Implementation {
	case "mysql":
		factory = new(mysql.Factory)
	default:
		err = errors.New("Invalid config, unknown param [db.implementation=" + config.DB.Implementation + "]")
	}
	return
}
