package main

import (
	"image"
	"sync"
	"time"

	"gocv.io/x/gocv"
)

type FrameData []byte

type Camera struct {
	vc        *gocv.VideoCapture
	output    *FanOut
	rec_cond  *sync.Cond
	framerate uint
}

func CreateCamera(framerate uint) (*Camera, error) {
	cam := new(Camera)
	camdev, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		return nil, err
	}

	cam.vc = camdev
	cam.output = CreateFanOut()
	cam.rec_cond = sync.NewCond(&sync.Mutex{})
	cam.framerate = framerate

	return cam, nil
}

func (cam *Camera) StartRecord() {
	go func() {
		img := gocv.NewMat()
		defer img.Close()
		delay := time.Second / time.Duration(cam.framerate)

		for {
			for !cam.output.Empty() {

				if ok := cam.vc.Read(&img); !ok {
					close(cam.output.input)
					return
				}
				if img.Empty() {
					continue
				}
				gocv.Resize(img, &img, image.Point{}, float64(0.5), float64(0.5), 0)
				frame, err := gocv.IMEncode(".jpg", img)

				if err == nil {
					cam.output.input <- frame.GetBytes()
					time.Sleep(delay)
				}
				frame.Close()
			}

			cam.rec_cond.L.Lock()
			cam.rec_cond.Wait()
		}
	}()
}

func (cam *Camera) FrameChan() chan FrameData {
	if cam.output.Empty() {
		defer cam.rec_cond.Signal()
	}
	return cam.output.Subscribe()
}

func (cam *Camera) ReleaseChan(ch chan FrameData) {
	cam.output.UnSubscribe(ch)
}
