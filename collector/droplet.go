package collector

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/prometheus/client_golang/prometheus"
)

type DropletCollector struct {
	client *godo.Client

	Up           *prometheus.Desc
	CPUs         *prometheus.Desc
	Memory       *prometheus.Desc
	Disk         *prometheus.Desc
	PriceHourly  *prometheus.Desc
	PriceMonthly *prometheus.Desc
}

func NewDropletCollector(client *godo.Client) *DropletCollector {
	labels := []string{"id", "name", "region"}

	return &DropletCollector{
		client: client,

		Up: prometheus.NewDesc(
			"digitalocean_droplet_up",
			"If 1 the droplet is up and running, 0 otherwise",
			labels, nil,
		),
		CPUs: prometheus.NewDesc(
			"digitalocean_droplet_cpus",
			"Droplet's number of CPUs",
			labels, nil,
		),
		Memory: prometheus.NewDesc(
			"digitalocean_droplet_memory_bytes",
			"Droplet's memory in bytes",
			labels, nil,
		),
		Disk: prometheus.NewDesc(
			"digitalocean_droplet_disk_bytes",
			"Droplet's disk in bytes",
			labels, nil,
		),
		PriceHourly: prometheus.NewDesc(
			"digitalocean_droplet_price_hourly",
			"Price of the Droplet billed hourly",
			labels, nil,
		),
		PriceMonthly: prometheus.NewDesc(
			"digitalocean_droplet_price_monthly",
			"Price of the Droplet billed monthly",
			labels, nil,
		),
	}
}

func (c *DropletCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Up
}

func (c *DropletCollector) Collect(ch chan<- prometheus.Metric) {
	droplets, _, err := c.client.Droplets.List(context.TODO(), nil)
	if err != nil {
		log.Printf("Can't list droplets: %v", err)
	}

	for _, droplet := range droplets {
		labels := []string{
			fmt.Sprintf("%d", droplet.ID),
			droplet.Name,
			droplet.Region.Slug,
		}

		var active float64
		if droplet.Status == "active" {
			active = 1.0
		}
		ch <- prometheus.MustNewConstMetric(
			c.Up,
			prometheus.GaugeValue,
			active,
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.CPUs,
			prometheus.GaugeValue,
			float64(droplet.Vcpus),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.Memory,
			prometheus.GaugeValue,
			float64(droplet.Memory*1024*1024),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.Disk,
			prometheus.GaugeValue,
			float64(droplet.Disk*1000*1000*1000),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.PriceHourly,
			prometheus.GaugeValue,
			float64(droplet.Size.PriceHourly), // TODO: Find out the correct data type
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.PriceMonthly,
			prometheus.GaugeValue,
			float64(droplet.Size.PriceMonthly), // TODO: Find out the correct data type
			labels...,
		)
	}
}