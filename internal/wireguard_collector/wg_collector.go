package wireguard_collector

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"golang.zx2c4.com/wireguard/wgctrl"
)

type wireguardCollector struct {
	logger log.Logger
	txDesc *prometheus.Desc
	rxDesc *prometheus.Desc
	lastDesc *prometheus.Desc
	friendlyNames map[string]string
}

func NewWireguardCollector(logger log.Logger, friendlyNames map[string]string) (prometheus.Collector, error) {
	return &wireguardCollector{
		logger: logger,
		txDesc: prometheus.NewDesc(
			prometheus.BuildFQName("wireguard", "", "sent_bytes_total"),
			"Bytes sent to the peer",
			[]string{"interface", "public_key", "friendly_name"},
			prometheus.Labels{},
		),
		rxDesc: prometheus.NewDesc(
			prometheus.BuildFQName("wireguard", "", "received_bytes_total"),
			"Bytes received from the peer",
			[]string{"interface", "public_key", "friendly_name"},
			prometheus.Labels{},
		),
		lastDesc: prometheus.NewDesc(
			prometheus.BuildFQName("wireguard", "", "latest_handshake_seconds"),
			"Seconds from the last handshake",
			[]string{"interface", "public_key", "friendly_name"},
			prometheus.Labels{},
		),
		friendlyNames: friendlyNames,
	}, nil
}

func (c *wireguardCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.rxDesc
	ch <- c.txDesc
	ch <- c.lastDesc
}

func (c *wireguardCollector) Collect(ch chan<- prometheus.Metric) {
	client, err := wgctrl.New()
	if err != nil {
		level.Error(c.logger).Log("msg", "Failed to create wireguard client object", "err", err.Error())
		return
	}
	defer client.Close()

	devices, err := client.Devices()
	if err != nil {
		level.Error(c.logger).Log("msg", "Failed to get wireguard devices", "err", err.Error())
		return
	}

	for _, device := range devices {
		for _, peer := range device.Peers {
			pubkey := peer.PublicKey.String()
			fname := c.friendlyNames[pubkey]

			ch <- prometheus.MustNewConstMetric(
				c.rxDesc,
				prometheus.CounterValue,
				float64(peer.ReceiveBytes),
				device.Name, pubkey, fname,
			)
			ch <- prometheus.MustNewConstMetric(
				c.txDesc,
				prometheus.CounterValue,
				float64(peer.TransmitBytes),
				device.Name, pubkey, fname,
			)
			ch <- prometheus.MustNewConstMetric(
				c.lastDesc,
				prometheus.GaugeValue,
				float64(peer.LastHandshakeTime.Unix()),
				device.Name, pubkey, fname,
			)
		}
	}
}
