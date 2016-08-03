package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/load"
)

var configFilepath string

func init() {
	flag.StringVar(&configFilepath, "c", "", "Config file to load")
}

func main() {
	flag.Parse()

	cfg := defaultConfig
	if configFilepath != "" {
		err := cfg.LoadFromFile(configFilepath)
		if err != nil {
			panic(err)
		}
	} else {
		panic("No configuration file specified")
	}

	stop := make(chan int)

	_, err := influxConnect("http://localhost:8086", stop)

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	go writeStats()

	// serialConfig := &serial.Config{
	// 	Name: cfg.SerialPort,
	// 	Baud: cfg.BaudRate,
	// }

	// serialPort, err := serial.OpenPort(serialConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// buf := make([]byte, 128)
	// for {
	// 	n, err := serialPort.Read(buf)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Printf("%q", buf[:n])
	// }

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	close(stop)
}

func writeStats() {
	// c := 8
	// b := make([]byte, c)
	for {
		// _, err := rand.Read(b)
		// if err != nil {
		// 	log.Fatalln("Error: ", err)
		// }
		avg, err := load.Avg()
		if err != nil {
			log.Fatalln("Error: ", err)
		}

		misc, err := load.Misc()
		if err != nil {
			log.Fatalln("Error: ", err)
		}

		rightNow := time.Now()

		tags := map[string]string{"machine": "robin"}

		loads := map[string]interface{}{
			"load1":  avg.Load1,
			"load5":  avg.Load5,
			"load15": avg.Load15,
		}
		createAndAddPoint("load", tags, loads, rightNow)

		miscs := map[string]interface{}{
			"ProcsRunning": misc.ProcsRunning,
			"ProcsBlocked": misc.ProcsBlocked,
			"Ctxt":         misc.Ctxt,
		}
		createAndAddPoint("misc", tags, miscs, rightNow)

		time.Sleep(1000 * time.Millisecond)
	}

	/*for {
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  "db",
			Precision: "ms",
		})

		if err != nil {
			log.Fatalln("Error: ", err)
		}

		// Create a point and add to batch
		tags := map[string]string{"cpu": "cpu-total"}
		fields := map[string]interface{}{
			"idle":   10.1,
			"system": 53.3,
			"user":   46.6,
		}
		pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())

		if err != nil {
			log.Fatalln("Error: ", err)
		}

		bp.AddPoint(pt)

		// Write the batch
		c.Write(bp)
		fmt.Println("stuff")
		time.Sleep(1 * time.Second)
	}*/
}
