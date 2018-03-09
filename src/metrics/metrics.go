package metrics

import (
	"fmt"
	"net/http"

	"../config"
	"../core"
	"../logging"
	"../stats/counters"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	serverCount             *prometheus.GaugeVec
	serverActiveConnections *prometheus.GaugeVec
	serverRxTotal           *prometheus.GaugeVec
	serverTxTotal           *prometheus.GaugeVec
	serverRxSecond          *prometheus.GaugeVec
	serverTxSecond          *prometheus.GaugeVec

	backendActiveConnections  *prometheus.GaugeVec
	backendRefusedConnections *prometheus.GaugeVec
	backendTotalConnections   *prometheus.GaugeVec
	backendRxBytes            *prometheus.GaugeVec
	backendTxBytes            *prometheus.GaugeVec
	backendRxSecond           *prometheus.GaugeVec
	backendTxSecond           *prometheus.GaugeVec
	backendLive               *prometheus.GaugeVec
)

func defineMetrics() {
	serverCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "server",
		Name:      "count",
		Help:      "Server Count.",
	}, []string{"server"})

	serverActiveConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "server",
		Name:      "active_connections",
		Help:      "Server Actice Connections.",
	}, []string{"server"})

	serverRxTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "server",
		Name:      "rx_total",
		Help:      "Server Rx Total.",
	}, []string{"server"})

	serverTxTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "server",
		Name:      "tx_total",
		Help:      "Server Tx Total.",
	}, []string{"server"})

	serverRxSecond = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "server",
		Name:      "rx_second",
		Help:      "Server Rx per Second.",
	}, []string{"server"})

	serverTxSecond = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "server",
		Name:      "tx_second",
		Help:      "Server Tx per Second.",
	}, []string{"server"})

	backendActiveConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "backend",
		Name:      "active_connections",
		Help:      "Backend Actice Connections.",
	}, []string{"server", "host", "port"})

	backendRefusedConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "backend",
		Name:      "refused_connections",
		Help:      "Backend Refused Connections.",
	}, []string{"server", "host", "port"})

	backendTotalConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "backend",
		Name:      "total_connections",
		Help:      "Backend Total Connections.",
	}, []string{"server", "host", "port"})

	backendRxBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "backend",
		Name:      "rx_bytes",
		Help:      "Backend Rx Bytes.",
	}, []string{"server", "host", "port"})

	backendTxBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "backend",
		Name:      "tx_bytes",
		Help:      "Backend Tx Bytes.",
	}, []string{"server", "host", "port"})

	backendRxSecond = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "backend",
		Name:      "rx_second",
		Help:      "Backend Rx per Second.",
	}, []string{"server", "host", "port"})

	backendTxSecond = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "backend",
		Name:      "tx_second",
		Help:      "Backend Tx per Second.",
	}, []string{"server", "host", "port"})

	backendLive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gobetween",
		Subsystem: "backend",
		Name:      "live",
		Help:      "Backend Alive.",
	}, []string{"server", "host", "port"})

}

func Start(cfg config.MetricsConfig) {

	var log = logging.For("metrics")

	if !cfg.Enabled {
		log.Info("Metrics disabled")
		return
	}

	log.Info("Starting up Metrics server ", cfg.Bind)
	defineMetrics()

	prometheus.MustRegister(serverCount)
	prometheus.MustRegister(serverActiveConnections)
	prometheus.MustRegister(serverRxTotal)
	prometheus.MustRegister(serverTxTotal)
	prometheus.MustRegister(serverRxSecond)
	prometheus.MustRegister(serverTxSecond)

	prometheus.MustRegister(backendActiveConnections)
	prometheus.MustRegister(backendRefusedConnections)
	prometheus.MustRegister(backendTotalConnections)
	prometheus.MustRegister(backendRxBytes)
	prometheus.MustRegister(backendTxBytes)
	prometheus.MustRegister(backendRxSecond)
	prometheus.MustRegister(backendTxSecond)
	prometheus.MustRegister(backendLive)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		fmt.Errorf("%s", http.ListenAndServe(cfg.Bind, nil))
	}()
}

func ReportHandleBackendLiveChange(server string, target core.Target, live bool) {
	if backendLive == nil {
		return
	}

	intLive := int(0)
	if live {
		intLive = 1
	}

	backendLive.WithLabelValues(server, target.Host, target.Port).Set(float64(intLive))
}

func ReportHandleConnectionsChange(server string, connections uint) {
	if serverActiveConnections == nil {
		return
	}

	serverActiveConnections.WithLabelValues(server).Set(float64(connections))
}

func ReportHandleStatsChange(server string, bs counters.BandwidthStats) {
	if serverRxTotal == nil || serverTxTotal == nil || serverRxSecond == nil || serverTxSecond == nil {
		return
	}

	serverRxTotal.WithLabelValues(server).Set(float64(bs.RxTotal))
	serverTxTotal.WithLabelValues(server).Set(float64(bs.TxTotal))
	serverRxSecond.WithLabelValues(server).Set(float64(bs.RxSecond))
	serverTxSecond.WithLabelValues(server).Set(float64(bs.TxSecond))
}

func ReportHandleBackendStatsChange(server string, target core.Target, backends map[core.Target]*core.Backend) {
	if serverCount == nil || backendRxBytes == nil || backendTxBytes == nil || backendRxSecond == nil || backendTxSecond == nil {
		return
	}

	backend, _ := backends[target]

	serverCount.WithLabelValues(server).Set(float64(len(backends)))

	backendRxBytes.WithLabelValues(server, target.Host, target.Port).Set(float64(backend.Stats.RxBytes))
	backendTxBytes.WithLabelValues(server, target.Host, target.Port).Set(float64(backend.Stats.TxBytes))
	backendRxSecond.WithLabelValues(server, target.Host, target.Port).Set(float64(backend.Stats.RxSecond))
	backendTxSecond.WithLabelValues(server, target.Host, target.Port).Set(float64(backend.Stats.TxSecond))
}

func ReportHandleOp(server string, target core.Target, backends map[core.Target]*core.Backend) {
	if backendActiveConnections == nil || backendRefusedConnections == nil || backendTotalConnections == nil {
		return
	}

	backend, _ := backends[target]

	backendActiveConnections.WithLabelValues(server, target.Host, target.Port).Set(float64(backend.Stats.ActiveConnections))
	backendRefusedConnections.WithLabelValues(server, target.Host, target.Port).Set(float64(backend.Stats.RefusedConnections))
	backendTotalConnections.WithLabelValues(server, target.Host, target.Port).Set(float64(backend.Stats.TotalConnections))
}
