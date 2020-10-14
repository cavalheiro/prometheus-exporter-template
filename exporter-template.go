package main

import (
	"flag"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// Config structure stores the values read from the TOML config
type (
	Config struct {
		BaseConfig configSection
	}
	configSection struct {
		Key1 uint64
		Key2 string
		Key3 []string
	}
)

var (
	config         Config
	addr           = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	debugFlag      = flag.Bool("debug", false, "Sets log level to debug.")
	configFileFlag = flag.String("config", "./config.toml", "Path to config file")
)

var (
	metric1 = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "metric_1_name",
		Help: "Matric 1 description",
	}, []string{"label1", "label2"})

	metric2 = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "metric_2_name",
			Help:    "Matric 2 description",
			Buckets: []float64{.1, .2, .3, .4, .5, .9, 1},
		},
	)
)

//Register metrics with Prometheus client
func init() {
	prometheus.MustRegister(metric1)
	prometheus.MustRegister(metric2)
}

// Load configuration
func init() {
	flag.Parse()
	// Setting logger to debug level when debug flag was set.
	if *debugFlag == true {
		log.SetLevel(log.DebugLevel)
	}
	// load config
	if _, err := os.Stat(*configFileFlag); err == nil {
		if _, err = toml.DecodeFile(*configFileFlag, &config); err != nil {
			log.Fatalf("Unable to parse configuration file: %s", err)
		}
	} else {
		log.Fatal("Please provide a config file with `--config <yourconfig>` or just create `config.toml` in this directory")
	}
	log.Infof("Configuration file settings: %+v", config)

}

// Periodically update metrics
func updateMetrics() {

	v := rand.Float64()

	metric1.WithLabelValues(
		"value1",
		"value2",
	).Set(v)

	metric2.Observe(v)

}

func main() {
	// Collect metrics on regular intervals
	go func() {
		for {
			updateMetrics()
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	// Serve via http
	http.Handle("/metrics", promhttp.Handler())
	log.Infof("Serving metrics on %s/metrics", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
