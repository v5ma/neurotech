package main

import (
	"fmt"

	"github.com/kataras/iris/websocket"
)

type WebsocketTunnel struct {
	cliconnections []websocket.Connection
}

func (wst *WebsocketTunnel) HandleEeg(c websocket.Connection) {
	fmt.Printf("websocket connection established with identifier: %s\n", c.ID())
	c.OnDisconnect(func() {
		fmt.Printf("websocket connection closed with identifer: %s\n", c.ID())
	})
	c.OnError(func(err error) {
		fmt.Printf("websocket connected error with identifier: %s\t%s\n", c.ID(), err)
	})
	c.OnMessage(func(data []byte) {
		for _, clic := range wst.cliconnections {
			clic.EmitMessage(data)
		}
	})
}

func (wst *WebsocketTunnel) HandleCli(c websocket.Connection) {
	fmt.Printf("websocket connection established with identifier: %s\n", c.ID())
	c.OnDisconnect(func() {
		fmt.Printf("websocket connection closed with identifer: %s\n", c.ID())
	})
	c.OnError(func(err error) {
		fmt.Printf("websocket connected error with identifier: %s\t%s\n", c.ID(), err)
	})
	wst.cliconnections = append(wst.cliconnections, c)
}
