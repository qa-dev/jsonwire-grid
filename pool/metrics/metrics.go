package metrics

import (
	log "github.com/Sirupsen/logrus"
	"jsonwire-grid/pool"
	"gopkg.in/alexcesaro/statsd.v2"
	"time"
)

type Sender struct {
	statd    *statsd.Client
	pool     *pool.Pool
	duration time.Duration
}

func NewSender(statd *statsd.Client, pool *pool.Pool, duration time.Duration) *Sender {
	return &Sender{statd, pool, duration}
}

func (s *Sender) SendAll() {
	for {
		s.countAvailableNodes()
		s.countTotalNodes()
		time.Sleep(s.duration)
	}
}

func (s *Sender) countTotalNodes() {
	count, err := s.pool.CountNodes(nil)
	if err != nil {
		log.Error("Can't get count total nodes: ", err.Error())
		return
	}
	s.statd.Gauge("node.total", count)
}

func (s *Sender) countAvailableNodes() {
	status := pool.NodeStatusAvailable
	count, err := s.pool.CountNodes(&status)
	if err != nil {
		log.Error("Can't get count total nodes: ", err.Error())
		return
	}
	s.statd.Gauge("node.available", count)
}
