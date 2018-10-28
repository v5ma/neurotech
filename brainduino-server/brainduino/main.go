package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/jacobsa/go-serial/serial"
)

const eegpath = "/ws/eeg"

var addr string
var brainduinopath string
var mock bool

func init() {
	flag.StringVar(&addr, "addr", "0.0.0.0:80", "addr to serve on")
	flag.StringVar(&brainduinopath, "brainduinopath", "/dev/rfcomm0", "path to brainduino serial device")
	flag.BoolVar(&mock, "mock", false, "to mock, or not to mock")
	flag.Parse()
}

func main() {
	// shutdown hook
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// init brainduino
	var device io.ReadWriteCloser
	var err error
	if !mock {
		device, err = serial.Open(serial.OpenOptions{
			PortName:              brainduinopath,
			BaudRate:              230400,
			InterCharacterTimeout: 100, // In milliseconds
			MinimumReadSize:       14,  // In bytes
			DataBits:              8,
			StopBits:              1,
		})
		if err != nil {
			fmt.Printf("Failed to open device: %s\n", err)
			return
		}
	} else {
		datastream := make(chan byte)
		go randomDatastream(datastream)
		device = mockDevice{
			datastream: datastream,
		}
	}

	b := NewBrainduino(device)

	rawlistener := make(chan interface{})
	b.Register(SampleListener, rawlistener)

	wsendpoint := url.URL{Scheme: "ws", Host: addr, Path: eegpath}
	fmt.Printf("Connecting to %s\n", wsendpoint.String())
	c, _, err := websocket.DefaultDialer.Dial(wsendpoint.String(), nil)
	if err != nil {
		fmt.Printf("error dial websocket eeg endpoint: %s\n", err)
		return
	}
	fmt.Printf("Connected to %s\n", wsendpoint.String())

	for {
		select {
		case d := <-rawlistener:
			d = d.(Sample)
			err = c.WriteJSON(d)
			if err != nil {
				fmt.Printf("error marshaling rawlistener data to json: %s\n", err)
			}
		case <-interrupt:
			fmt.Println("Shutting down")
			b.Unregister(SampleListener, rawlistener)
			b.Close()
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
				return
			}
			c.Close()
			return
		}
	}
}
