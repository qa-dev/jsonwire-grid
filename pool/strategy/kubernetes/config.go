package kubernetes

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/qa-dev/jsonwire-grid/config"
)

type strategyParams struct {
	Namespace          string
	PodCreationTimeout time.Duration
}

func (sp *strategyParams) UnmarshalJSON(b []byte) error {
	tempStruct := struct {
		Namespace          string `json:"namespace"`
		PodCreationTimeout string `json:"pod_creation_timeout"`
	}{
		"default",
		"1m",
	}
	if err := json.Unmarshal(b, &tempStruct); err != nil {
		return err
	}
	podCreationTimeout, err := time.ParseDuration(tempStruct.PodCreationTimeout)
	if err != nil {
		return fmt.Errorf("invalid value strategy.pod_creation_timeout in config, given: %v", tempStruct.PodCreationTimeout)
	}
	sp.Namespace = tempStruct.Namespace
	sp.PodCreationTimeout = podCreationTimeout
	return nil
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
