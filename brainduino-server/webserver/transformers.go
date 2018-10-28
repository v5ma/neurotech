package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"time"

	pubsub "github.com/dustin/go-broadcast"
	"github.com/mjibson/go-dsp/fft"
)

type ListenerType int

const (
	SampleListener ListenerType = iota
	FFTListener
)

type Subscriber interface {
	Register(ListenerType, chan<- interface{})
	Unregister(ListenerType, chan<- interface{})
}

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

func (b *Brainduino) fftloop() {
	// assumes b.numchan == 2
	ctr := 0
	var seqnum uint
	fftsize := 256
	fftdata0 := make([]float64, fftsize)
	fftdata1 := make([]float64, fftsize)
	rawlistener := make(chan interface{})
	b.rawBroadcaster.Register(rawlistener)
	for {
		s := <-rawlistener
		sample := s.(Sample)
		fftdata0[ctr%fftsize] = sample.Channels[0]
		fftdata1[ctr%fftsize] = sample.Channels[1]
		// Set the frequency that the FFT is sent out.
		// e.g. ctr%2==0, every other sample
		//      ctr%10==0, every 10th sample
		//      ctr%250==0, every 250th sample
		if ctr%4 == 0 {
			fftd := FFTData{
				Name:           "fft",
				Channels:       make([][]float64, b.numchan),
				SequenceNumber: seqnum,
				Timestamp:      time.Now(),
			}
			fftd.Channels[0] = abs(fft.FFTReal(fftdata0))[:125]
			fftd.Channels[1] = abs(fft.FFTReal(fftdata1))[:125]
			b.fftBroadcaster.Submit(fftd)
			seqnum++
		}
		ctr++

	}
}

func (b Brainduino) Register(t ListenerType, listener chan<- interface{}) {
	switch t {
	case SampleListener:
		b.rawBroadcaster.Register(listener)
	case FFTListener:
		b.fftBroadcaster.Register(listener)
	}
}

func (b Brainduino) Unregister(t ListenerType, listener chan<- interface{}) {
	switch t {
	case SampleListener:
		b.rawBroadcaster.Unregister(listener)
	case FFTListener:
		b.fftBroadcaster.Unregister(listener)
	}
}
