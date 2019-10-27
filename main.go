package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/qa-dev/jsonwire-grid/config"
	"github.com/qa-dev/jsonwire-grid/handlers"
	"github.com/qa-dev/jsonwire-grid/logger"
	"github.com/qa-dev/jsonwire-grid/middleware"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	poolMetrics "github.com/qa-dev/jsonwire-grid/pool/metrics"
	"github.com/qa-dev/jsonwire-grid/utils/metrics"
)

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	cfg := config.New()
	err := cfg.LoadFromFile(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("Problem in loading config from file, %s", err)
	}
	err = logger.Init(cfg.Logger)
	if err != nil {
		log.Fatalf("Problem in init logger, %s", err)
	}

	middlewareWrap := middleware.NewWrap(log.StandardLogger())

	busyNodeDuration, err := time.ParseDuration(cfg.Grid.BusyNodeDuration)
	if err != nil {
		log.Fatal("Invalid value grid.busy_node_duration in config")
	}
	reservedNodeDuration, err := time.ParseDuration(cfg.Grid.ReservedDuration)
	if err != nil {
		log.Fatal("Invalid value grid.reserved_node_duration in config")
	}
	storageFactory, err := invokeStorageFactory(*cfg)
	if err != nil {
		log.Fatalf("Can't invoke storage factory, %s", err)
	}
	storage, err := storageFactory.Create(*cfg)
	if err != nil {
		log.Fatalf("Can't create storage factory, %s", err)
	}
	capsComparator := capabilities.NewComparator()
	strategyFactoryList, err := invokeStrategyFactoryList(*cfg)
	if err != nil {
		log.Fatalf("Can't create strategy factory list, %s", err)
	}

	clientFactory, err := createClient(*cfg)
	if err != nil {
		log.Fatalf("Create ClientFactory error, %s", err)
	}

	var strategyList []pool.StrategyInterface
	for _, strategyFactory := range strategyFactoryList {
		strategy, err := strategyFactory.Create(storage, capsComparator, clientFactory)
		if err != nil {
			log.Fatalf("Can't create strategy, %s", err)
		}
		strategyList = append(strategyList, strategy)
	}

	strategyListStruct := pool.NewStrategyList(strategyList)
	poolInstance := pool.NewPool(storage, strategyListStruct)
	poolInstance.SetBusyNodeDuration(busyNodeDuration)
	poolInstance.SetReservedNodeDuration(reservedNodeDuration)

	go func() {
		for {
			poolInstance.FixNodeStatuses()
			time.Sleep(time.Minute * 5) // todo: move to config
		}
	}()

	cache := pool.NewCache(time.Minute * 10) // todo: move to config

	go func() {
		for {
			cache.CleanUp()
			time.Sleep(time.Minute) // todo: move to config
		}
	}()

	if cfg.Statsd != nil {
		statsdClient, err := metrics.NewStatsd(
			cfg.Statsd.Host,
			cfg.Statsd.Port,
			cfg.Statsd.Protocol,
			cfg.Statsd.Prefix,
			cfg.Statsd.Enable)
		// todo: move to config
		poolMetricsSender := poolMetrics.NewSender(statsdClient, poolInstance, time.Second*1, cfg.Statsd.CapabilitiesList, capsComparator)
		go poolMetricsSender.SendAll()
		if err != nil {
			log.Errorf("Statsd create socked error: %s", err)
		}
		middlewareWrap.Add(middleware.NewStatsd(log.StandardLogger(), statsdClient, true).RegisterMetrics)
	}

	http.Handle("/wd/hub/session", middlewareWrap.Do(&handlers.CreateSession{Pool: poolInstance, ClientFactory: clientFactory})) //selenium
	http.Handle("/session", middlewareWrap.Do(&handlers.CreateSession{Pool: poolInstance, ClientFactory: clientFactory}))        //wda
	http.Handle("/grid/register", middlewareWrap.Do(&handlers.RegisterNode{Pool: poolInstance}))
	http.Handle("/grid/status", middlewareWrap.Do(&handlers.GridStatus{Pool: poolInstance, Config: *cfg}))
	http.Handle("/grid/api/proxy", &handlers.APIProxy{Pool: poolInstance})
	http.HandleFunc("/_info", heartbeat)
	http.Handle("/", middlewareWrap.Do(&handlers.UseSession{Pool: poolInstance, Cache: cache}))

	server := &http.Server{Addr: fmt.Sprintf(":%v", cfg.Grid.Port)}
	serverError := make(chan error)
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			serverError <- err
		}
	}()

	select {
	case err = <-serverError:
		log.Fatalf("Server error, %s", err)
	case <-stop:
	}

	log.Info("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // todo: move to config
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("graceful shutdown, %v", err)
	}

	log.Info("Server gracefully stopped")
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"result": {"ok": true}}`))
	if err != nil {
		log.Errorf("write response, %v", err)
	}
}
