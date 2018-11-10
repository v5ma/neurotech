package main

import (
	"encoding/json"
	"time"

	pubsub "github.com/dustin/go-broadcast"
	"github.com/kataras/iris/websocket"
	"github.com/mjibson/go-dsp/fft"
)

type Sample struct {
	Name           string
	Channels       []float64
	Timestamp      time.Time
	SequenceNumber uint
}

type FFTData struct {
	Name           string
	Channels       [][]float64
	Timestamp      time.Time
	SequenceNumber uint
}

type WebsocketTunnel struct {
	pubsub.Broadcaster
}

func NewWebsocketTunnel() *WebsocketTunnel {
	return &WebsocketTunnel{
		Broadcaster: pubsub.NewBroadcaster(256 * 256),
	}
}

func (wst *WebsocketTunnel) HandleEeg(c websocket.Connection) {
	LOG.Infof("websocket connection established with identifier: %s", c.ID())
	rawlistener := make(chan []byte, 256*256)
	donech := make(chan bool)
	go wst.fftloop(rawlistener, donech)
	c.OnDisconnect(func() {
		donech <- true
		LOG.Infof("websocket connection closed with identifer: %s", c.ID())
	})
	c.OnError(func(err error) {
		LOG.Errorf("websocket connection error with identifier: %s\t[%s]", c.ID(), err)
	})
	c.OnMessage(func(data []byte) {
		rawlistener <- data
	})
}

func (wst *WebsocketTunnel) HandleCli(c websocket.Connection) {
	LOG.Infof("websocket connection established with identifier: %s", c.ID())
	fftch := make(chan interface{})
	donech := make(chan bool)
	wst.Register(fftch)
	c.OnDisconnect(func() {
		wst.Unregister(fftch)
		donech <- true
		LOG.Infof("websocket connection closed with identifier: %s", c.ID())
	})
	c.OnError(func(err error) {
		LOG.Errorf("websocket connected error with identifier: %s\t%s", c.ID(), err)
	})
	go clisten(c, fftch, donech)
}

func clisten(c websocket.Connection, fftch chan interface{}, donech chan bool) {
	LOG.Infof("Started client listener on websocket connection with identifier: %s", c.ID())
	for {
		select {
		case d := <-fftch:
			c.EmitMessage(d.([]byte))
		case <-donech:
			LOG.Infof("Closing client listener on websocket connection with identifier: %s", c.ID())
			return
		}
	}
}

func (wst *WebsocketTunnel) fftloop(rawlistener chan []byte, donech chan bool) {
	// assumes b.numchan == 2
	numchan := 2
	ctr := 0
	var seqnum uint
	fftsize := 256
	fftdata0 := make([]float64, fftsize)
	fftdata1 := make([]float64, fftsize)
	for {
		select {
		case s := <-rawlistener:
			sample := &Sample{}
			err := json.Unmarshal(s, sample)
			if err != nil {
				LOG.Errorf("error unmarshalling sample: %s", err)
			}
			fftdata0[ctr%fftsize] = sample.Channels[0]
			fftdata1[ctr%fftsize] = sample.Channels[1]
			// Set the frequency that the FFT is sent out.
			// e.g. ctr%2==0, every other sample
			//      ctr%10==0, every 10th sample
			//      ctr%250==0, every 250th sample
			if ctr%16 == 0 {
				fftd := FFTData{
					Name:           "fft",
					Channels:       make([][]float64, numchan),
					SequenceNumber: seqnum,
					Timestamp:      time.Now(),
				}
				fftd.Channels[0] = abs(fft.FFTReal(fftdata0))[:125]
				fftd.Channels[1] = abs(fft.FFTReal(fftdata1))[:125]
				jsonfft, err := json.Marshal(fftd)
				if err != nil {
					LOG.Errorf("error marshalling fft data: %s\n", err)
					continue
				}
				wst.Submit(jsonfft)
				seqnum++
			}
			ctr++
		case <-donech:
			LOG.Info("Closing fftloop")
			return
		}
	}
}
