package main

import (
	"encoding/json"
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

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
