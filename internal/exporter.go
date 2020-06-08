package exporter

import (
	"log"
	"net/http"
	"time"

	//"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Exporter struct {
	//remaining_date prometheus.Gauge
	SSLMetrics map[string]*prometheus.Desc
	Config     *BaseConfig
}

func (s *Server) Start() {
	log.Fatal(http.ListenAndServe(":8081", s.Handler))
}

func NewServer(exporter Exporter) *Server {
	log.Println(`
 	 This is a prometheus exporter for stream
  	Access: http://127.0.0.1:8081
  	`)
	r := http.NewServeMux()
	metricsPath := "/metrics"
	prometheus.MustRegister(&exporter)
	r.Handle(metricsPath, promhttp.Handler())
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		 <head><title>Dummy Exporter</title></head>
		 <body>
		 <h1>Stream Exporter</h1>
		 <p><a href='` + metricsPath + `'>Metrics</a></p>
		 </body>
		 </html>`))
	})

	return &Server{Handler: r, exporter: exporter}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	log.Println("Collect()")

	data, err := e.gatherData()
	log.Println(data)
	if err != nil {
		log.Fatalf("Error for openssl: %v", err)
		//log.Errorf("Error gathering Data from Mysql server: %v", err)
		return
	}

	err = e.processMetrics(data, ch)
	if err != nil {
		log.Fatalf("Error Processing Metrics", err)
	}
	log.Println("All Metrics successfully collected.")
	/*
		e.remaining_date.Set(float64(10))
		e.remaining_date.Collect(ch)
	*/
}

// 讓exporter的prometheus屬性呼叫Describe方法

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	//e.remaining_date.Describe(ch)
	//e.gauge_metrics2.Describe(ch)
	for _, m := range e.SSLMetrics {
		ch <- m
	}
}

func (e *Exporter) processMetrics(data SSLInfoArray, ch chan<- prometheus.Metric) error {
	log.Println("processMetrics()")
	local1, _ := time.LoadLocation("Asia/Taipei")
	for _, x := range data {
		//增加label的地方
		ch <- prometheus.MustNewConstMetric(e.SSLMetrics["certificate_remaining_date"], prometheus.GaugeValue, x.SSLRemainingDate, x.DomainNmme, x.ExpiredDate.In(local1).Format("2006-01-02 15:04:05"), x.RegistryDate.In(local1).Format("2006-01-02 15:04:05"))

	}
	return nil

}

func AddMetrics() map[string]*prometheus.Desc {

	SSLMetrics := make(map[string]*prometheus.Desc)

	SSLMetrics["certificate_remaining_date"] = prometheus.NewDesc(
		prometheus.BuildFQName("certificate", "", "remaining_date"),
		"A metric with a constant '0' value labeled by matchId,roomId,locationId from DB table.",
		[]string{"domainName", "expiredDate", "registryDate"}, nil,
	)

	log.Println("Metrics added.....")

	return SSLMetrics
}
