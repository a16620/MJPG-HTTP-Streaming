package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	cam, err := CreateCamera(15)

	if err != nil {
		log.Fatalln("Camera setting failed")
	}

	cam.StartRecord()

	host := "0.0.0.0:8091"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		frame := cam.FrameChan()
		defer cam.ReleaseChan(frame)

		w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
		for {
			data := "--frame\r\n  Content-Type: image/jpeg\r\n\r\n" + string(<-frame) + "\r\n\r\n"
			w.Write([]byte(data))
		}
	})

	fmt.Println("start")
	log.Fatal(http.ListenAndServe(host, nil))
}
