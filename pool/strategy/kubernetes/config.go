package kubernetes

import (
	"encoding/json"
	"errors"
	"github.com/qa-dev/jsonwire-grid/config"
)

type strategyParams struct {
}

type strategyConfig struct {
	config.Strategy
	Params   strategyParams
	NodeList []nodeConfig
}

type nodeConfig struct {
	config.Node
	Params nodeParams
}

type nodeParams struct {
	Image string `json:"image"`
	Port  string `json:"port"`
}

func newConfig(cfg config.Strategy) (*strategyConfig, error) {
	kubConfig := new(strategyConfig)
	kubConfig.Type = cfg.Type
	kubConfig.Limit = cfg.Limit
	err := json.Unmarshal(cfg.Params, &kubConfig.Params)
	if err != nil {
		return nil, errors.New("unmarshal strategy params, " + err.Error())
	}
	kubConfig.NodeList = make([]nodeConfig, len(cfg.NodeList))
	for i, defaultCfgNode := range cfg.NodeList {
		kubConfig.NodeList[i].CapabilitiesList = cfg.NodeList[i].CapabilitiesList
		err = json.Unmarshal(defaultCfgNode.Params, &kubConfig.NodeList[i].Params)
		if err != nil {
			return nil, errors.New("unmarshal node params, " + err.Error())
		}
	}
	return kubConfig, nil
}
