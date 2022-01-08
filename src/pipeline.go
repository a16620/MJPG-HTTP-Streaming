package main

import (
	"sync"
	"time"
)

type ns struct{} //null struct

type FanOut struct {
	sync.RWMutex
	subscribers map[chan FrameData]ns
	input       chan FrameData
}

func CreateFanOut() *FanOut {
	fanout := new(FanOut)

	fanout.subscribers = make(map[chan FrameData]ns)
	fanout.input = make(chan FrameData, 2)

	go func() {
		for frame := range fanout.input {
			fanout.RLock()

			for subs := range fanout.subscribers {
				select {
				case subs <- frame:
				case <-time.After(time.Millisecond * 10):
				}
			}

			fanout.RUnlock()
		}
	}()

	return fanout
}

func (fanout *FanOut) Subscribe() chan FrameData {
	fanout.Lock()
	defer fanout.Unlock()

	newChan := make(chan FrameData, 2)
	fanout.subscribers[newChan] = ns{}

	return newChan
}

func (fanout *FanOut) UnSubscribe(channel chan FrameData) {
	fanout.Lock()
	defer fanout.Unlock()

	delete(fanout.subscribers, channel)
}

func (fanout *FanOut) Empty() bool {
	fanout.RLock()
	defer fanout.RUnlock()

	return len(fanout.subscribers) == 0
}
