package kubernetes

import (
	"errors"
	"github.com/qa-dev/jsonwire-grid/config"
	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type StrategyFactory struct {
	Config config.Strategy
}

func (f *StrategyFactory) Create(
	storage pool.StorageInterface,
	capsComparator capabilities.ComparatorInterface,
	clientFactory jsonwire.ClientFactoryInterface,
) (pool.StrategyInterface, error) {
	strategyConfig, err := newConfig(f.Config)
	if err != nil {
		return nil, errors.New("convert strategy config to k8s format, " + err.Error())
	}

	//todo: При применении конфига возможно стоит проверять запушены ли образы
	for _, nodeCfg := range f.Config.NodeList {
		for _, capsCfg := range nodeCfg.CapabilitiesList {
			capsComparator.Register(capsCfg)
		}
	}

	//todo: выпилить этот говноклиент, когда будет работать нормальный
	kubConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.New("create k8s config, " + err.Error())
	}

	clientset, err := kubernetes.NewForConfig(kubConfig)
	if err != nil {
		return nil, errors.New("create k8s clientset, " + err.Error())
	}

	provider := &kubDnsProvider{
		clientset:     clientset,
		namespace:     "default", //todo: брать из конфига !!!
		clientFactory: clientFactory,
	}

	return &Strategy{
		storage,
		provider,
		*strategyConfig,
		capsComparator,
	}, nil
}
