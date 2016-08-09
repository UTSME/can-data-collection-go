package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UTSME/go-serial/serial"
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

	// port, err := serial.Open(options)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// defer port.Close()
	//
	// portScanner := bufio.NewScanner(port)
	//
	// for portScanner.Scan() {
	// 	c := &CANMessage{}
	// 	c.parse(portScanner.Text())
	// 	log.Println(c)
	// }

	// for {
	// 	line, err := portReader.ReadString(delim)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	n := len(line)
	// 	log.Printf("%q", line[:n])
	// }

}
