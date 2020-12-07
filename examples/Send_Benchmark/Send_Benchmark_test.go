package main

import (
	"../.."
	"testing"
)

func BenchmarkSendVideoFrame(b *testing.B) {
	sender := NDI.SendCreate(&NDI.SendCreateT{
		PNdiName:   "benchmark",
		ClockAudio: false,
		ClockVideo: false,
	})
	if sender == nil {
		panic("unable to create sender")
	}

	defer func() {
		NDI.SendSendVideoAsyncV2(sender, nil)
		NDI.SendDestroy(sender)
		NDI.Destroy()
	}()

	src := createData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transmitFrame(sender, &src[xRes*2*(i%(yRes*(scrollDist-1)))])
	}
}
