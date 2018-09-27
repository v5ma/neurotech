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
	/*
		device := mockDevice{
			datastream: make(chan byte),
		}
		go randomDatastream(device.datastream)
	*/
	b := NewBrainduino(device)
	defer b.Close()

	rawlistener := make(chan interface{})
	b.Register(SampleListener, rawlistener)
	defer b.Unregister(SampleListener, rawlistener)

	fftlistener := make(chan interface{})
	b.Register(FFTListener, fftlistener)
	defer b.Unregister(FFTListener, fftlistener)

	app := iris.New()

	// set up http routes
	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile(indexfile, false)
	})
	app.Post("/command/{id:string}", func(ctx iris.Context) {
		id := ctx.Params().Get("id")
		if !isValidCommand(id) {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Writef("Internal Server Error: %s command not supported\n", id)
			return
		}
		_, err := b.Write([]byte(id))
		if err != nil {
			fmt.Printf("Error writing to brainduino: %s\n", err)
		}
	})

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
