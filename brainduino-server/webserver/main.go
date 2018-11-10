package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/websocket"
)

var addr string
var indexfile string
var logfile string
var chartsngraphsfile string

var LOG *golog.Logger

func init() {
	flag.StringVar(&addr, "addr", "0.0.0.0:80", "url to serve on")
	flag.StringVar(&indexfile, "indexfile", "./static/index.html", "path to index.html")
	flag.StringVar(&logfile, "logfile", "/var/log/brainduino/webserver.log", "path to webserver.log")
	flag.StringVar(&chartsngraphsfile, "chartsngraphsfile", "./static/chartsngraphs.html", "path to chartsngraphs.html")
	flag.Parse()
}

func main() {
	app := iris.New()

	// set up logging
	err := os.MkdirAll(path.Dir(logfile), 0755)
	if err != nil {
		fmt.Printf("Error making directory %s for logfile: %s\n", path.Dir(logfile), err)
	}
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening logfile %s: %s\n", logfile, err)
	}

	defaultLogger := app.Logger()
	defaultLogger.SetOutput(f)
	LOG = defaultLogger
	defaultLogger.Info("server starting")

	customLogger := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
		// Query appends the url query to the Path.
		Query: true,

		//Columns: true,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKeys: []string{"logger_message"},

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		MessageHeaderKeys: []string{"User-Agent"},
	})
	app.Use(customLogger)

	// set up http routes
	app.StaticWeb("/static", "./static")
	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile(indexfile, false)
	})
	app.Get("/chartsngraphs", func(ctx iris.Context) {
		ctx.ServeFile(chartsngraphsfile, false)
	})

	// set up websocket routes
	wst := NewWebsocketTunnel()
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

	defaultLogger.Info("Server started")
	app.Run(iris.Addr(addr))
}
