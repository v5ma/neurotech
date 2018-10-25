package main

import (
	"flag"

	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

var url string
var indexfile string
var chartsngraphsfile string

func init() {
	flag.StringVar(&url, "url", "0.0.0.0:80", "url to serve on")
	flag.StringVar(&indexfile, "indexfile", "../static/index.html", "path to index.html")
	flag.StringVar(&chartsngraphsfile, "chartsngraphsfile", "./static/chartsngraphs.html", "path to chartsngraphs.html")
	flag.Parse()
}

func main() {
	app := iris.New()

	// set up http routes
	app.StaticWeb("/static", "./static")
	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile(indexfile, false)
	})
	app.Get("/chartsngraphs", func(ctx iris.Context) {
		ctx.ServeFile(chartsngraphsfile, false)
	})
	/*
		app.Post("/command/{id:string}", func(ctx iris.Context) {
			id := ctx.Params().Get("id")
			if !isValidCommand(id) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Writef("Internal Server Error: %s command not supported\n", id)
				return
			}
			_, err := b.Write([]byte(id))
			if err != nil {
				fmt.printf("Error writing to brainduino: %s\n", err)
			}
		})
	*/

	wst := WebsocketTunnel{
		cliconnections: make([]websocket.Connection, 0),
	}

	// set up websocket routes
	wseeg := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
	wseeg.OnConnection(wst.HandleEeg) // handle client connections

	wscli := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
	wscli.OnConnection(wst.HandleCli) // handle client connections

	app.Get("/ws/eeg", wseeg.Handler())
	app.Get("/ws/client", wscli.Handler())

	app.Run(iris.Addr(url))
}
