package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/cmplx"

	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

var url string
var indexfile string
var brainduinopath string

func rootHandler(ctx iris.Context) {
	ctx.ServeFile(indexfile, false)
}

type WebsocketTunnel struct {
	rawlistener <-chan interface{}
	fftlistener <-chan interface{}
	connections []websocket.Connection
}

func (wst *WebsocketTunnel) Handle(c websocket.Connection) {
	fmt.Printf("websocket connection established with identifier: %s\n", c.ID())
	c.OnDisconnect(func() {
		fmt.Printf("websocket connection closed with identifer: %s\n", c.ID())
	})
	c.OnError(func(err error) {
		fmt.Printf("websocket connected error with identifier: %s\t%s\n", c.ID(), err)
	})
	wst.connections = append(wst.connections, c)
}

func abs(cin []complex128) []float64 {
	fout := make([]float64, len(cin))
	for idx, v := range cin {
		fout[idx] = cmplx.Abs(v)
	}
	return fout
}

func (wst *WebsocketTunnel) broadcast() {
	for {
		select {
		case s := <-wst.rawlistener:
			d := s.(Sample)
			for _, c := range wst.connections {
				jsondata, err := json.Marshal(d)
				if err != nil {
					fmt.Printf("error marshalling sample json: %s\n", err)
					continue
				}
				c.EmitMessage(jsondata)
			}
		case f := <-wst.fftlistener:
			d := f.(FFTData)
			for _, c := range wst.connections {
				jsondata, err := json.Marshal(d)
				if err != nil {
					fmt.Printf("error marshalling sample json: %s\n", err)
					continue
				}
				c.EmitMessage(jsondata)
			}
		}
	}
}

func init() {
	flag.StringVar(&url, "url", "0.0.0.0:8080", "url to serve on")
	flag.StringVar(&indexfile, "indexfile", "./static/index.html", "path to index.html")
	flag.StringVar(&brainduinopath, "brainduinopath", "/dev/rfcomm0", "path to brainduino serial device")
	flag.Parse()
}

func main() {
	// init brainduino
	/*
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
	*/
	device := mockDevice{
		datastream: make(chan byte),
	}
	b := NewBrainduino(device)
	go randomDatastream(device.datastream)
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
