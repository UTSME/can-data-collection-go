package main

import (
	"bufio"
	"flag"
	"log"

	"github.com/UTSME/go-serial/serial"
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

	options := serial.OpenOptions{
		PortName:        cfg.SerialPort,
		BaudRate:        cfg.BaudRate,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	port, err := serial.Open(options)
	if err != nil {
		log.Fatal(err)
	}

	defer port.Close()

	portScanner := bufio.NewScanner(port)

	for portScanner.Scan() {
		can, _ := parseCANMessage(portScanner.Text())
		log.Println(can)
	}

	// for {
	// 	line, err := portReader.ReadString(delim)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	n := len(line)
	// 	log.Printf("%q", line[:n])
	// }

}
