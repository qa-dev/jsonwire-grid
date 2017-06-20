package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/qa-dev/jsonwire-grid/config"

	"errors"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/strategy/persistent"
	"github.com/qa-dev/jsonwire-grid/storage/mysql"
)

type StorageFactoryInterface interface {
	Create(config.Config) (pool.StorageInterface, error)
}

type StrategyFactoryInterface interface {
	Create(config.Config, pool.StorageInterface) (pool.StrategyInterface, error)
}

func invokeStorageFactory(config config.Config) (factory StorageFactoryInterface, err error) {
	switch config.DB.Implementation {
	case "mysql":
		factory = new(mysql.Factory)
	default:
		err = errors.New("Invalid config, unknown param [db.implementation=" + config.DB.Implementation + "]")
	}
	return
}

func invokeStrategyFactoryList(config config.Config, storage pool.StorageInterface) (factoryList []StrategyFactoryInterface, err error) {
	for _, strategyConfig := range config.Grid.StrategyList {
		switch strategyConfig.Type {
		case string(pool.NodeTypePersistent):
			factoryList = append(factoryList, new(persistent.StrategyFactory))
		//case string(pool.NodeTypeKubernetes):
			//factoryList = append(factoryList, new(kubernetes.StrategyFactory))
		default:
			err = errors.New("Undefined strategy type: " + strategyConfig.Type)
			return
		}
	}
	return
}
