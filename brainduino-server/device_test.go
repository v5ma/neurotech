package main

import (
	"testing"
	"time"
)

func TestOffsetBinaryToInt(t *testing.T) {
	tables := []struct {
		hexstr   []byte
		expected int
	}{
		{[]byte("000000"), -8388608}, // all off
		{[]byte("FFFFFF"), 8388607},  // all on
		{[]byte("800000"), 0},        // 1000...0000
		{[]byte("800001"), 1},        // 1000...0001
		{[]byte("7FFFFF"), -1},       // 0111...0000
	}

	datastream := make(chan byte)
	b, _ := newMockBrainduino(datastream)

	for _, table := range tables {
		actual := b.offsetBinaryToInt(table.hexstr)
		if actual != table.expected {
			t.Errorf("For: %s, Got: %d, Want: %d\n", table.hexstr, actual, table.expected)
		}
	}
}

func TestADCNorm(t *testing.T) {
	tables := []struct {
		raw      int
		expected float64
	}{
		{-8388608, -5.0},
		{8388607, 5.0},
		{-4194304, -2.5},
		{4194304, 2.5},
		{0, 0.0},
	}

	datastream := make(chan byte)
	b, _ := newMockBrainduino(datastream)

	for _, table := range tables {
		actual := b.adcnorm(table.raw)
		if actual < table.expected-0.0001 || actual > table.expected+0.0001 {
			t.Errorf("For: %d, Got: %f, Want: %f\n", table.raw, actual, table.expected)
		}
	}
}

func TestReadloop(t *testing.T) {
	tables := []struct {
		testdata []byte
		expected []Sample
	}{
		{[]byte("000000\t000000\r000000\t000000\r000000\t000000\r"), []Sample{Sample{[]float64{-5.0, -5.0}, time.Now(), 0}, Sample{[]float64{-5.0, -5.0}, time.Now(), 0}, Sample{[]float64{-5.0, -5.0}, time.Now(), 0}}},
	}

	for _, table := range tables {
		datastream := make(chan byte)
		mbd, _ := newMockBrainduino(datastream)
		testlistener := make(chan Sample)
		mbd.RegisterListener("testlistener", testlistener)
		go func() {
			for _, d := range table.testdata {
				datastream <- d
			}
			defer mbd.Close()
		}()
		for actual := range testlistener {
			for _, sample := range table.expected {
				for channum, channel := range sample.Channels {
					if actual.Channels[channum] != channel {
						t.Errorf("For: %x, Got: %f, Want: %f\n", table.testdata, actual.Channels[channum], channel)
					}
				}
			}
			close(testlistener)
		}
	}
}
