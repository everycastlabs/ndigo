package main

import (
	"fmt"
	"strings"
	"time"

	NDI "github.com/broadcastervc/ndigo"
)

var latencies = make([]float32, 0, 100)

func receiveSource(sourceName *NDI.SourceT, referenceTime time.Time) {
	recv := NDI.RecvCreateV3(nil)
	NDI.RecvConnect(recv, sourceName)

	for {
		var videoFrame NDI.VideoFrameV2T
		if NDI.RecvCaptureV2(recv, &videoFrame, nil, nil, 250) == NDI.FrameTypeVideo {
			videoFrame.Deref()
			sinceReference := time.Since(referenceTime)

			// Should we exit?
			if strings.Contains(videoFrame.PMetadata, "<exit/>") {
				break
			}

			latencyMs := float32(sinceReference.Microseconds()-videoFrame.Timecode) / 1000
			latencies = append(latencies, latencyMs)

			fmt.Printf("NDI video latency (with compression, transmission and decompression) = %1.2fms\n", latencyMs)
		}
		videoFrame.Free()
	}

	NDI.RecvDestroy(recv)
}

func main() {
	if !NDI.Initialize() {
		panic("cannot run NDI")
	}

	sender := NDI.SendCreate(&NDI.SendCreateT{
		PNdiName:   "Video and Audio",
		ClockVideo: false,
	})
	if sender == nil {
		panic("unable to create sender")
	}

	referenceTime := time.Now()
	// start the receiver
	go receiveSource(NDI.SendGetSourceName(sender), referenceTime)

	videoData := make([]byte, 640*480*2)
	videoFrame := &NDI.VideoFrameV2T{
		Xres:   640,
		Yres:   480,
		FourCC: NDI.FourCCTypeUYVY,
		PData:  &videoData[0],
	}

	ticker := time.NewTicker(time.Millisecond * 250)

	for idx := 100; idx > 0; idx-- {
		videoFrame.Timecode = time.Since(referenceTime).Microseconds()
		if idx == 1 {
			videoFrame.PMetadata = "<exit/>"
		}

		NDI.SendSendVideoV2(sender, videoFrame)

		videoFrame.Free()

		<-ticker.C
	}

	NDI.SendSendVideoV2(sender, nil)
	NDI.SendDestroy(sender)
	NDI.Destroy()

	var averageLatency float32
	for _, l := range latencies {
		averageLatency += l
	}
	averageLatency = averageLatency / float32(len(latencies))

	fmt.Printf("Average NDI video latency = %1.2fms", averageLatency)
}
