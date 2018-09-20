package main

import (
	"flag"
	"fmt"

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
	listeners []chan Sample
}

func (wst *WebsocketTunnel) Handle(c websocket.Connection) {
	fmt.Printf("websocket connection established with identifier: %s\n", c.ID())
	c.OnDisconnect(func() {
		fmt.Printf("websocket connection closed with identifer: %s\n", c.ID())
	})
	c.OnError(func(err error) {
		fmt.Printf("websocket connected error with identifier: %s\t%s\n", c.ID(), err)
	})
}

func (wst *WebsocketTunnel) broadcast(c websocket.Connection) {
	for {
		for _, listener := range wst.listeners {
			select {
			case <-listener:
				c.To(websocket.Broadcast).EmitMessage([]byte{'\x00'})
			}
		}

	}
}

func init() {
	flag.StringVar(&url, "url", "0.0.0.0:8080", "url to serve on")
	flag.StringVar(&indexfile, "indexfile", "./static/webxr1.html", "path to index.html")
	flag.StringVar(&brainduinopath, "brainduinopath", "/dev/rfcomm0", "path to brainduino serial device")
	flag.Parse()
}

func main() {
	// init brainduino
	b, err := NewBrainduino(brainduinopath)
	if err != nil {
		fmt.Errorf("Failed to open brainduino: %s\n", err)
	}
	defer b.Close()

	timeserieslistener := make(chan Sample)
	b.RegisterListener("timeseries_ws_listener", timeserieslistener)

	fftlistener := make(chan Sample)
	b.RegisterListener("fft_ws_listener", fftlistener)

	app := iris.New()

	// set up http routes
	app.Get("/", rootHandler)

	// set up websocket routes
	wst := &WebsocketTunnel{
		listeners: []chan Sample{timeserieslistener, fftlistener},
	}
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
	ws.OnConnection(wst.Handle)
	app.Get("/ws", ws.Handler())

	app.Run(iris.Addr(url))
}
