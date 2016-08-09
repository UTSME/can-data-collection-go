package main

import (
	"bytes"
	"log"

	"github.com/UTSME/go-cobs"
)

func scanFrame(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, 0x00); i >= 0 {
		// We have a full null-terminated frame.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated frame. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func decodeFrames(frames chan []byte, stop chan int) chan []byte {
	packets := make(chan []byte, 32)
	go func() {
		for {
			select {
			case <-stop:
				close(packets)
				return
			case frame := <-frames:
				packet, err := cobs.Decode(frame)
				if err != nil {
					log.Fatal(err)
					continue
				}
				packets <- packet
			}
		}
	}()
	return packets
}
