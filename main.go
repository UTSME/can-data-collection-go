package main

import (
	"flag"
	"log"

	"github.com/UTSME/serial"
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

	serialConfig := &serial.Config{
		Name: cfg.SerialPort,
		Baud: cfg.BaudRate,
	}

	serialPort, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 128)
	for {
		n, err := serialPort.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%q", buf[:n])
	}

}
