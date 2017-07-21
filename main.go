package main

import (
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/jsonwire-grid/config"
	"github.com/qa-dev/jsonwire-grid/handlers"
	"github.com/qa-dev/jsonwire-grid/logger"
	"github.com/qa-dev/jsonwire-grid/middleware"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	poolMetrics "github.com/qa-dev/jsonwire-grid/pool/metrics"
	"github.com/qa-dev/jsonwire-grid/utils/metrics"
	"net/http"
	"os"
	"os/signal"
	"time"
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

	statsdClient, err := metrics.NewStatsd(
		cfg.Statsd.Host,
		cfg.Statsd.Port,
		cfg.Statsd.Protocol,
		cfg.Statsd.Prefix,
		cfg.Statsd.Enable)

	if nil != err {
		log.Errorf("Statsd create socked error: %s", err)
	}

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

	poolMetricsSender := poolMetrics.NewSender(statsdClient, poolInstance, time.Second*1) // todo: move to config
	go poolMetricsSender.SendAll()

	go func() {
		for {
			poolInstance.FixNodeStatuses()
			time.Sleep(time.Minute * 5) // todo: move to config
		}
	}()

	m := middleware.NewLogMiddleware(statsdClient)
	http.Handle("/wd/hub/session", m.Log(&handlers.CreateSession{Pool: poolInstance, ClientFactory: clientFactory})) //selenium
	http.Handle("/session", m.Log(&handlers.CreateSession{Pool: poolInstance, ClientFactory: clientFactory}))        //wda
	http.Handle("/grid/register", m.Log(&handlers.RegisterNode{Pool: poolInstance}))
	http.Handle("/grid/api/proxy", &handlers.APIProxy{Pool: poolInstance})
	http.HandleFunc("/_info", heartbeat)
	http.Handle("/", m.Log(&handlers.UseSession{Pool: poolInstance}))

	server := &http.Server{Addr: fmt.Sprintf(":%v", cfg.Grid.Port)}
	serverError := make(chan error)
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			// todo: норма ли что при вызове server.Shutdown всегда возвращается еррор???
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
