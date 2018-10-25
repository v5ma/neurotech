package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jacobsa/go-serial/serial"
)

var url string
var chartsngraphsfile string
var brainduinopath string
var mock bool

func init() {
	flag.StringVar(&url, "url", "0.0.0.0:80", "url to serve on")
	flag.StringVar(&brainduinopath, "brainduinopath", "", "path to brainduino serial device")
	flag.StringVar(&chartsngraphsfile, "chartsngraphsfile", "./static/chartsngraphs.html", "path to chartsngraphs.html")
	flag.BoolVar(&mock, "mock", false, "to mock, or not to mock")
	flag.Parse()
}

func getSystemBrainduinoDevicePath() string {
	if len(brainduinopath) > 0 {
		return brainduinopath
	} else {
		basestr := "/dev/rfcomm"
		for i := 0; i < 10; i++ {
			basestr = strings.Join([]string{basestr, strconv.Itoa(i)}, "")
			// On Debian, this will not work as when the device is connected, disconnected and reconnected there are now 2
			// devices, /dev/rfcomm0 and /dev/rfcomm1. Perhaps in the FileInfo we can decipher which is the brainduino rfcomm device.
			if _, err := os.Stat(basestr); !os.IsNotExist(err) {
				return basestr
			}
		}
	}
	return "not found"
}

func main() {
	// init brainduino
	var device io.ReadWriteCloser
	var err error
	if !mock {
		device, err = serial.Open(serial.OpenOptions{
			PortName:              getSystemBrainduinoDevicePath(),
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
	defer b.Close()

	rawlistener := make(chan interface{})
	b.Register(SampleListener, rawlistener)
	defer b.Unregister(SampleListener, rawlistener)

	fftlistener := make(chan interface{})
	b.Register(FFTListener, fftlistener)
	defer b.Unregister(FFTListener, fftlistener)
	
	for {
		switch {
		case d := <- rawlistener:
		}
	}
}
