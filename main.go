package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UTSME/go-serial/serial"
	"github.com/shirou/gopsutil/load"
)

var configFilepath string

func init() {
	flag.StringVar(&configFilepath, "c", "", "Config file to load")
}

func main() {
	flag.Parse()

	stop := make(chan int)

	cfg := defaultConfig
	if configFilepath != "" {
		err := cfg.LoadFromFile(configFilepath)
		if err != nil {
			panic(err)
		}
	} else {
		panic("No configuration file specified")
	}

	_, err := influxConnect("http://localhost:8086", stop)

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	go writeStats()

	options := serial.OpenOptions{
		PortName:        cfg.SerialPort,
		BaudRate:        cfg.BaudRate,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	var cans chan *CANMessage

	if false {
		frames := scanSerial(options, stop)
		packets := decodeFrames(frames, stop)
		cans = parseMessages(packets, stop)
	} else {
		cans = generateCanMessages(stop, time.Second)
	}

	go func(_cans chan *CANMessage) {
		for {
			log.Println(<-_cans)
		}
	}(cans)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	close(stop)
}

func writeStats() {
	for {
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
}
