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
	logger.Init(cfg.Logger)

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
	reservedNodeDuration, err := time.ParseDuration(cfg.Grid.BusyNodeDuration)
	if err != nil {
		log.Fatal("Invalid value grid.busy_node_duration in config")
	}
	storage, err := invokeStorage(*cfg)
	if err != nil {
		log.Fatalf("Can't invoke storage, %s", err)
	}
	poolInstance := pool.NewPool(storage)
	poolInstance.SetBusyNodeDuration(busyNodeDuration)
	poolInstance.SetReservedNodeDuration(reservedNodeDuration)

	//todo: сделать конфиг для пула, вынести duration туда
	poolMetricsSender := poolMetrics.NewSender(statsdClient, poolInstance, time.Second*1) //todo: вынести в конфиг
	go poolMetricsSender.SendAll()

	go func() {
		for {
			poolInstance.FixNodeStatuses()
			time.Sleep(time.Minute * 5) // todo: вынести в конфиг
		}
	}()

	m := middleware.NewLogMiddleware(statsdClient)
	http.Handle("/wd/hub/session", m.Log(&handlers.CreateSession{Pool: poolInstance})) //selenium
	http.Handle("/session", m.Log(&handlers.CreateSession{Pool: poolInstance}))        //wda
	http.Handle("/grid/register", m.Log(&handlers.RegisterNode{Pool: poolInstance}))
	http.Handle("/grid/api/proxy", &handlers.ApiProxy{Pool: poolInstance})
	http.HandleFunc("/_info", heartbeat)
	http.Handle("/", m.Log(&handlers.UseSession{Pool: poolInstance}))

	server := &http.Server{Addr: fmt.Sprintf(":%v", cfg.Grid.Port)}
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			// todo: норма ли что при закрытии всегда возвращается еррор???
			log.Errorf("Listen serve error, %s", err)
		}
	}()

	<-stop

	log.Info("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // todo: вынести в конфиг
	defer cancel()
	server.Shutdown(ctx)

	log.Info("Server gracefully stopped")
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"result": {"ok": true}}`))
}
