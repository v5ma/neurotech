package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/jacobsa/go-serial/serial"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

var url string
var indexfile string
var chartsngraphsfile string
var brainduinopath string
var mock bool

func init() {
	flag.StringVar(&url, "url", "0.0.0.0:80", "url to serve on")
	flag.StringVar(&indexfile, "indexfile", "./static/index.html", "path to index.html")
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

	app := iris.New()

	// set up http routes
	app.StaticWeb("/static", "./static")
	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile(indexfile, false)
	})
	app.Get("/chartsngraphs", func(ctx iris.Context) {
		ctx.ServeFile(chartsngraphsfile, false)
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
