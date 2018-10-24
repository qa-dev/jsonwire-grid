package metrics

import (
	log "github.com/sirupsen/logrus"
	"github.com/qa-dev/jsonwire-grid/pool"
	"github.com/qa-dev/jsonwire-grid/pool/capabilities"
	"gopkg.in/alexcesaro/statsd.v2"
	"time"
)

// Sender - metrics sender.
type Sender struct {
	statd          *statsd.Client
	pool           *pool.Pool
	duration       time.Duration
	selectorList   []CapabilitiesSelector
	capsComparator capabilities.ComparatorInterface
}

type CapabilitiesSelector struct {
	Tag          string                    `json:"tag"`
	Capabilities capabilities.Capabilities `json:"capabilities"`
}

// NewSender - constructor of sender.
func NewSender(
	statd *statsd.Client,
	pool *pool.Pool,
	duration time.Duration,
	capsMetricList []CapabilitiesSelector,
	capsComparator capabilities.ComparatorInterface,
) *Sender {
	return &Sender{statd, pool, duration, capsMetricList, capsComparator}
}

// NewSender - sends metrics of nodes availability.
func (s *Sender) SendAll() {
	for {
		s.countAvailableNodes()
		s.countTotalNodes()
		s.countByCapabilities()
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

func (s *Sender) countByCapabilities() {
	nodeList, err := s.pool.GetAll()
	if err != nil {
		log.Error("Can't get all nodes: ", err.Error())
		return
	}

	for _, node := range nodeList {
		for _, availableCaps := range node.CapabilitiesList {
			s.capsComparator.Register(availableCaps)
		}
	}

	for _, requiredCaps := range s.selectorList {
		availableCount := 0
		reservedCount := 0
		busyCount := 0
		for _, node := range nodeList {
			for _, availableCaps := range node.CapabilitiesList {
				if s.capsComparator.Compare(requiredCaps.Capabilities, availableCaps) {
					switch node.Status {
					case pool.NodeStatusAvailable:
						availableCount++
					case pool.NodeStatusReserved:
						reservedCount++
					case pool.NodeStatusBusy:
						busyCount++
					}
					break
				}
			}
		}
		s.statd.Gauge("node-by-caps."+requiredCaps.Tag+".available", availableCount)
		s.statd.Gauge("node-by-caps."+requiredCaps.Tag+".reserved", reservedCount)
		s.statd.Gauge("node-by-caps."+requiredCaps.Tag+".busy", busyCount)
	}
}
