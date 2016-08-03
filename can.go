package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type CANMessage struct {
	Identifier uint32
	Timestamp  time.Time
	Data       []uint8
}

func parseCANMessage(s string) (*CANMessage, error) {

	s1 := strings.Split(s, ":")
	if len(s1) != 2 {
		return nil, errors.New("Invalid number of :")
	}

	id64, err := strconv.ParseUint(s1[0], 0, 32)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		dataArray[i] = uint8(data64)
	}

	can := &CANMessage{
		Identifier: identifier,
		Timestamp:  time.Now(),
		Data:       dataArray,
	}
	return can, nil
}
