package main

import (
	"bufio"
	"log"

	"github.com/UTSME/go-serial/serial"
)

func scanSerial(opts serial.OpenOptions, stop chan int) chan []byte {
	port, err := serial.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	frames := make(chan []byte, 32)

	scanner := bufio.NewScanner(port)
	scanner.Split(scanFrame)

	go func() {
		<-stop
		port.Close()
		close(frames)
	}()

	//the meaty part
	go func() {
		for scanner.Scan() {
			frames <- scanner.Bytes()
		}
	}()

	return frames
}
