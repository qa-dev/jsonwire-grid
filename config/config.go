package config

import (
	"encoding/json"
	"errors"
	"github.com/Sirupsen/logrus"
	log "github.com/Sirupsen/logrus"
	"os"
)

type Config struct {
	Logger Logger `json:"logger"`
	DB     DB     `json:"db"`
	Grid   Grid   `json:"grid"`
	Statsd Statsd `json:"statsd"`
}

type Grid struct {
	Port             int        `json:"port"`
	StrategyList     []Strategy `json:"strategy_list"`
	BusyNodeDuration string     `json:"busy_node_duration"`     // duration string format ex. 12m, see time.ParseDuration()
	ReservedDuration string     `json:"reserved_node_duration"` // duration string format ex. 12m, see time.ParseDuration()
}

type Strategy struct {
	Config   map[string]string `json:"config"` // ex. docker config, kubernetes config, etc.
	Type     string            `json:"type"`
	Limit    int               `json:"limit"`
	NodeList []Node            `json:"node_list"`
}

type Node struct {
	Config       map[string]string      `json:"config"` // ex. image_name, etc.
	Capabilities map[string]interface{} `json:"capabilities"`
	Limit        int                    `json:"limit"`
}

type Logger struct {
	Level logrus.Level `json:"level"`
}

type DB struct {
	Implementation string `json:"implementation"`
	Connection     string `json:"connection"`
}

type Statsd struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Prefix   string `json:"prefix"`
	Enable   bool   `json:"enable"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) LoadFromFile(path string) error {
	log.Printf(path)
	if path == "" {
		return errors.New("empty configuration file path")
	}

	configFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)

	return jsonParser.Decode(&c)
}
