package metrics

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"
)

// NewStatsd возвращает новый настроенный клиент statsd
func NewStatsd(host string, port int, protocol string, prefix string, enable bool) (*statsd.Client, error) {
	protocol = strings.ToLower(protocol)
	muted := !enable

	log.Infof(
		`Create statsd client to %v:%v via %v with prefix "%v", muted is %v.`,
		host, port, protocol, prefix, muted)

	client, err := statsd.New(
		statsd.Address(fmt.Sprintf("%v:%v", host, port)),
		statsd.Prefix(prefix),
		statsd.Network(protocol),
		statsd.Mute(muted))

	log.Info("Statsd client was created.")

	if err != nil {
		return nil, err
	}

	return client, nil
}
