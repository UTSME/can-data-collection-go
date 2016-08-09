package main

import (
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
	"time"
)

type CANMessage struct {
	Identifier uint32
	Timestamp  time.Time
	Data       []byte
}

func parseMessages(packets chan []byte, stop chan int) chan *CANMessage {
	cans := make(chan *CANMessage, 32)
	go func() {
		for {
			select {
			case <-stop:
				close(cans)
				return
			case packet := <-packets:
				c := &CANMessage{}
				c.parsePacket(packet)
				cans <- c
			}
		}
	}()
	return cans
}

func (c *CANMessage) parsePacket(packet []byte) error {
	if len(packet) != 12 {
		return errors.New("Wrong frame length")
	}
	c.Identifier = binary.LittleEndian.Uint32(packet[0:3])
	c.Timestamp = time.Now()
	c.Data = packet[4:11]
	return nil
}

func (c *CANMessage) parseString(s string) error {

	s1 := strings.Split(s, ":")
	if len(s1) != 2 {
		return errors.New("Invalid number of :")
	}

	id64, err := strconv.ParseUint(s1[0], 0, 32)
	if err != nil {
		return err
	}

	identifier := uint32(id64)

	allData := strings.Split(s1[1], ",")
	if len(allData) > 8 {
		allData = allData[:8]
	}

	dataArray := make([]uint8, len(allData))

	for i, oneData := range allData {
		data64, err := strconv.ParseUint(oneData, 0, 8)
		if err != nil {
			return err
		}
		dataArray[i] = uint8(data64)
	}

	c.Identifier = identifier
	c.Timestamp = time.Now()
	c.Data = dataArray
	return nil
}

func generateCanMessages(stop chan int, d time.Duration) chan *CANMessage {
	cans := make(chan *CANMessage)
	go func() {
		for {
			select {
			case <-stop:
				close(cans)
				return
			default:
				cans <- &CANMessage{
					Identifier: 1200,
					Timestamp:  time.Now(),
					Data:       []uint8{0, 1, 2, 3, 4, 5, 6, 7},
				}
				time.Sleep(d)
			}
		}
	}()
	return cans
}
