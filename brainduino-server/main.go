package main

import (
	"flag"
	"fmt"

	"github.com/jacobsa/go-serial/serial"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

var url string
var indexfile string
var brainduinopath string

func init() {
	flag.StringVar(&url, "url", "0.0.0.0:8080", "url to serve on")
	flag.StringVar(&indexfile, "indexfile", "./static/index.html", "path to index.html")
	flag.StringVar(&brainduinopath, "brainduinopath", "/dev/rfcomm0", "path to brainduino serial device")
	flag.Parse()
}

func main() {
	// init brainduino
	device, err := serial.Open(serial.OpenOptions{
		PortName:              "/dev/rfcomm0",
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
	//device := mockDevice{
	//	datastream: make(chan byte),
	//}
	b := NewBrainduino(device)
	//go randomDatastream(device.datastream)
	defer b.Close()

	rawlistener := make(chan interface{})
	b.RegisterRawListener(rawlistener)
	fftlistener := make(chan interface{})
	b.RegisterFFTListener(fftlistener)

	app := iris.New()

	// set up http routes
	app.Get("/", rootHandler)

	// set up websocket routes
	wst := &WebsocketTunnel{
		rawlistener: rawlistener,
		fftlistener: fftlistener,
		connections: make([]websocket.Connection, 0),
	}
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
	ws.OnConnection(wst.Handle)
	go wst.broadcast()
	app.Get("/ws", ws.Handler())

	app.Run(iris.Addr(url))
}
