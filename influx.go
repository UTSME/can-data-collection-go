package main

import (
	"fmt"
	"log"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

const (
	writeCacheTimeout = 10
)

var points chan *influx.Point
var client influx.Client

func influxConnect(addr string, stop chan int) (influx.Client, error) {
	var err error
	points = make(chan *influx.Point)
	client, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr: addr,
	})
	go startBatchWrite(stop)
	return client, err
}

func createAndAddPoint(name string,
	tags map[string]string,
	fields map[string]interface{},
	t time.Time) {

	p, err := influx.NewPoint(name, tags, fields, t)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	addPointToQueue(p)
}

func addPointToQueue(point *influx.Point) {
	points <- point
}

func startBatchWrite(stop chan int) {
	timeout := make(chan int)
	for {
		fmt.Println("new batch")
		batch, err := influx.NewBatchPoints(influx.BatchPointsConfig{
			Database:  "db",
			Precision: "ms",
		})
		go func() {
			time.Sleep(time.Second * writeCacheTimeout)
			timeout <- 0
		}()

		if err != nil {
			log.Fatalln("Error: ", err)
		}

		fmt.Println("top")

		func() {
			for {
				select {
				case <-stop:
					writeBatch(batch)
					return
				case point := <-points:
					batch.AddPoint(point)
					continue
				case <-timeout:
					writeBatch(batch)
					return
				}
			}
		}()
	}
}

func writeBatch(batch influx.BatchPoints) {
	fmt.Println("writing", len(batch.Points()))
	client.Write(batch)
}
