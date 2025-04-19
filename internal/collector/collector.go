package collector

import (
	"log"

	"github.com/604pierce/ratgdo-exporter/internal/client"
	"github.com/prometheus/client_golang/prometheus"
)

type MyCollector struct {
	apiBaseURL  string
	deviceLabel string

	// These hold the metric descriptors
	statusUp        *prometheus.Desc
	upTime          *prometheus.Desc
	openingsCount   *prometheus.Desc
	freeHeap        *prometheus.Desc
	garageDoorState *prometheus.Desc // Optional: can expose as a label if needed
	garageLockState *prometheus.Desc
	openDuration    *prometheus.Desc
	garageLightOn   *prometheus.Desc
}

func (c *MyCollector) Collect(ch chan<- prometheus.Metric) {
	ratgdoClient := client.NewClient(c.apiBaseURL)
	if err := ratgdoClient.Authenticate(); err != nil {
		log.Printf("Auth failed: %v", err)
		ch <- prometheus.MustNewConstMetric(c.statusUp, prometheus.GaugeValue, 0)
		return
	}

	status, err := ratgdoClient.GetStatus()
	if err != nil {
		log.Printf("Failed to get status: %v", err)
		ch <- prometheus.MustNewConstMetric(c.statusUp, prometheus.GaugeValue, 0, "unknown")
		return
	}
	log.Printf("Response status: %v", status)

	ch <- prometheus.MustNewConstMetric(c.statusUp, prometheus.GaugeValue, 1, status.DeviceName)
	ch <- prometheus.MustNewConstMetric(c.upTime, prometheus.GaugeValue, status.UpTime, status.DeviceName)
	ch <- prometheus.MustNewConstMetric(c.openingsCount, prometheus.GaugeValue, status.OpeningsCount, status.DeviceName)
	ch <- prometheus.MustNewConstMetric(c.freeHeap, prometheus.GaugeValue, status.FreeHeap, status.DeviceName)
	ch <- prometheus.MustNewConstMetric(c.openDuration, prometheus.GaugeValue, status.OpenDuration, status.DeviceName)
	ch <- prometheus.MustNewConstMetric(c.garageDoorState, prometheus.GaugeValue, doorStateToValue(status.GarageDoorState), status.DeviceName)
	ch <- prometheus.MustNewConstMetric(c.garageLightOn, prometheus.GaugeValue, lightStateToValue(status.GarageLight), status.DeviceName)
}

func NewCollector(apiBaseURL string) *MyCollector {
	return &MyCollector{
		apiBaseURL:  apiBaseURL,
		deviceLabel: "Hello World",
		statusUp: prometheus.NewDesc(
			"ratgdo_up",
			"Whether the exporter could reach the ratgdo API",
			[]string{"device"}, nil,
		),
		upTime: prometheus.NewDesc(
			"ratgdo_uptime_seconds",
			"Uptime of the device",
			[]string{"device"}, nil,
		),
		openingsCount: prometheus.NewDesc(
			"ratgdo_openings_total",
			"Number of garage door openings",
			[]string{"device"}, nil,
		),
		freeHeap: prometheus.NewDesc(
			"ratgdo_free_heap_bytes",
			"Amount of free heap memory on the device",
			[]string{"device"}, nil,
		),
		garageDoorState: prometheus.NewDesc(
			"ratgdo_garage_door_state",
			"State of the garage door, either open, closed, or unknown",
			[]string{"device"}, nil,
		),
		garageLockState: prometheus.NewDesc(
			"ratgdo_garage_door_lock_state",
			"Describes if the garage door is locked or not",
			[]string{"device"}, nil,
		),
		openDuration: prometheus.NewDesc(
			"ratgdo_open_duration",
			"Duration of a lift cycle",
			[]string{"device"}, nil,
		),
		garageLightOn: prometheus.NewDesc(
			"ratgdo_garage_light_on",
			"Status of the garage light",
			[]string{"device"}, nil,
		),
	}
}

func (c *MyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.statusUp
	ch <- c.upTime
	ch <- c.openingsCount
	ch <- c.freeHeap
	ch <- c.garageDoorState
	ch <- c.garageLockState
	ch <- c.openDuration
	ch <- c.garageLightOn
}

func doorStateToValue(state string) float64 {
	switch state {
	case "Open":
		return 1
	case "Closed":
		return 0
	default:
		return -1
	}
}

func lightStateToValue(state bool) float64 {
	if state {
		return 1
	} else {
		return 0
	}
}
