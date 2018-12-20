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

var (
	addr    string
	addrint string
	addrext string
	logfile string
	LOG     *golog.Logger
)

func init() {
	flag.StringVar(&addr, "addr", "0.0.0.0:80", "internal url to serve on")
	flag.StringVar(&addrint, "addrint", "", "internal url to serve on")
	flag.StringVar(&addrext, "addrext", "", "external url to serve on")
	flag.StringVar(&logfile, "logfile", "", "path to log file, if no path then stdout")
	flag.Parse()
	if addrint == "" {
		addrint = addr
	}
	if addrext == "" {
		addrext = addr
	}
}

func main() {
	app := iris.New()

	// set up logging
	err := os.MkdirAll(path.Dir(logfile), 0755)
	if err != nil {
		fmt.Printf("Error making directory %s for logfile: %s\n", path.Dir(logfile), err)
	}
	var logf *os.File
	if logfile == "" {
		logf = os.Stdout
	} else {
		logf, err = os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening logfile %s: %s\n", logfile, err)
		}
	}

	defaultLogger := app.Logger()
	defaultLogger.SetOutput(logf)
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

		// will be added to the logs.
		MessageContextKeys: []string{"logger_message"},

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		MessageHeaderKeys: []string{"User-Agent"},
	})
	app.Use(customLogger)

	// set up static http routes
	app.StaticWeb("/static/imgs", "./static/imgs")
	app.RegisterView(iris.HTML("./static/views", ".html").Reload(true))
	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("addr", addrext)
		ctx.View("index.html")
	})
	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("addr", addrext)
		ctx.View("chartsngraphs.html")
	})

	// set up rest routes
	deviceRegistrations := NewDeviceRegistrations()
	app.Post("/device/registration", PostRegistration(deviceRegistrations))
	app.Get("/device/registration/{id:string}", GetRegistration(deviceRegistrations))

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

	defaultLogger.Info("server started")
	app.Run(iris.Addr(addrint), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:   true,
		EnableOptimizations: true,
	}))
}
