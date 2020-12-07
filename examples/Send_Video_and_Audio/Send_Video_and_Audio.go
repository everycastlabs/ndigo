package main

import (
	"../.."
	"fmt"
	"math/rand"
	"time"
	"unsafe"
)

func main() {
	if !NDI.Initialize() {
		panic("cannot run NDI")
	}

	sender := NDI.SendCreate(&NDI.SendCreateT{
		PNdiName: "Video and Audio",
	})
	if sender == nil {
		panic("unable to create sender")
	}

	// We are going to create a 1920x1080 interlaced frame at 29.97Hz.
	videoData := make([]byte, 1920*1080*2)
	videoFrame := &NDI.VideoFrameV2T{
		Xres:   1920,
		Yres:   1080,
		FourCC: NDI.FourCCTypeUYVY,
		PData:  &videoData[0],
	}

	// Because 48kHz audio actually involves 1601.6 samples per frame, we make a basic sequence that we follow.
	audioNoSamples := []int{1602, 1601, 1602, 1601, 1602}

	audioData := make([]float32, 1602*2)
	audioFrame := &NDI.AudioFrameV2T{
		SampleRate:           48000,
		NoChannels:           2,
		NoSamples:            1602,
		PData:                &audioData[0],
		ChannelStrideInBytes: int32(unsafe.Sizeof(float32(0)) * 1602),
	}

	defer func() {
		videoFrame.Free()
		audioFrame.Free()
		NDI.SendDestroy(sender)
		NDI.Destroy()
	}()

	ticker := time.NewTicker(time.Second / 50)

	for idx := 0; idx < 1000; idx++ {
		<-ticker.C
		// display black
		black := (idx % 50) > 10

		// Because we are clocking to the video it is better to always submit the audio
		// before, although there is very little in it. I'll leave it as an excercies for the
		// reader to work out why.
		audioFrame.NoSamples = int32(audioNoSamples[idx%5])

		// When not black, insert noise into the buffer. This is a horrible noise, but its just
		// for illustration.
		// Fill in the buffer with silence. It is likely that you would do something much smarter than this.
		for as := range audioData {
			if black {
				audioData[as] = 0
			} else {
				audioData[as] = rand.Float32() * 2
			}
		}

		// Submit the audio buffer
		NDI.SendSendAudioV2(sender, audioFrame)

		// Every 50 frames display a few frames of white
		for vs := 0; vs < len(videoData); vs += 2 {
			var v uint16
			if black {
				v = 128 | (16 << 8)
			} else {
				v = 128 | (235 << 8)
			}
			videoData[vs] = byte(v)
			videoData[vs+1] = byte(v >> 8)
		}

		NDI.SendSendVideoV2(sender, videoFrame)
		fmt.Printf("Frame number %d sent\n", idx)
	}
}
