package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/qa-dev/jsonwire-grid/config"

	"errors"
	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"github.com/qa-dev/jsonwire-grid/pool/strategy/kubernetes"
	"github.com/qa-dev/jsonwire-grid/pool/strategy/persistent"
	"github.com/qa-dev/jsonwire-grid/selenium"
	"github.com/qa-dev/jsonwire-grid/storage/local"
	"github.com/qa-dev/jsonwire-grid/storage/mysql"
	"github.com/qa-dev/jsonwire-grid/wda"
)

type StorageFactoryInterface interface {
	Create(config.Config) (pool.StorageInterface, error)
}

type StrategyFactoryInterface interface {
	Create(pool.StorageInterface, capabilities.ComparatorInterface, jsonwire.ClientFactoryInterface) (pool.StrategyInterface, error)
}

func invokeStorageFactory(config config.Config) (factory StorageFactoryInterface, err error) {
	switch config.DB.Implementation {
	case "mysql":
		factory = new(mysql.Factory)
	case "local":
		factory = new(local.Factory)
	default:
		err = errors.New("Invalid config, unknown param [db.implementation=" + config.DB.Implementation + "]")
	}
	return
}

func invokeStrategyFactoryList(config config.Config) (factoryList []StrategyFactoryInterface, err error) {
	for _, strategyConfig := range config.Grid.StrategyList {
		switch strategyConfig.Type {
		case string(pool.NodeTypePersistent):
			factoryList = append(factoryList, new(persistent.StrategyFactory))
		case string(pool.NodeTypeKubernetes):
			factoryList = append(factoryList, &kubernetes.StrategyFactory{Config: strategyConfig})
		default:
			err = errors.New("Undefined strategy type: " + strategyConfig.Type)
			return
		}
	}
	return
}

func createClient(config config.Config) (jsonwire.ClientFactoryInterface, error) {
	switch config.Grid.ClientType {
	case "wda":
		return new(wda.ClientFactory), nil
	case "selenium":
		return new(selenium.ClientFactory), nil
	default:
		return nil, errors.New("undefined value config.Grid.ClientType")
	}
}
