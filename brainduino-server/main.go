package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"

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
	listeners   []chan Sample
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
	// We may not always want to broadcast Samples???
	fftdata := make([]float64, 250)
	ctr := 0
	for {
		for _, listener := range wst.listeners {
			select {
			case sample := <-listener:
				fftdata[ctr%250] = sample.Channels[0]
				sample.FFT = abs(fft.FFTReal(fftdata))[:124]
				for _, c := range wst.connections {
					samplejson, err := json.Marshal(sample)
					if err != nil {
						fmt.Printf("error marshalling sample json: %s\n", err)
						continue
					}
					c.To(websocket.Broadcast).EmitMessage(samplejson)
				}
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
	b, err := NewBrainduino(brainduinopath)
	if err != nil {
		fmt.Printf("Failed to open brainduino: %s\n", err)
		return
	}
	defer b.Close()

	timeserieslistener := make(chan Sample)
	b.RegisterListener("timeseries_ws_listener", timeserieslistener)

	app := iris.New()

	// set up http routes
	app.Get("/", rootHandler)

	// set up websocket routes
	wst := &WebsocketTunnel{
		listeners:   []chan Sample{timeserieslistener},
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
