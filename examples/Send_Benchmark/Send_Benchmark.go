package main

import (
	"fmt"
	"math"
	"runtime/debug"
	"time"

	NDI "github.com/broadcastervc/ndigo"
)

import "C"

const (
	xRes       = 3840
	yRes       = 2160
	framerateN = 60000
	framerateD = 1001
	scrollDist = 4
)

func clamp(x float64) int {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 255
	}
	return int(x * 255)
}

func createData() []uint8 {
	src := make([]uint8, xRes*yRes*scrollDist*2)
	for y := 0.0; y < yRes*scrollDist; y++ {
		for x := 0.0; x < xRes; x += 2 {
			i := int(y*xRes*2 + x*2)

			// generate some patterns of some kind
			fy := y / yRes
			fx_0 := (x + 0) / xRes
			fx_1 := (x + 1) / xRes

			// get the RGB colour
			r0 := clamp(math.Cos(fx_0*9.0+fy*9.5)*0.5 + 0.5)
			g0 := clamp(math.Cos(fx_0*12.0+fy*40.5)*0.5 + 0.5)
			b0 := clamp(math.Cos(fx_0*23.0+fy*15.5)*0.5 + 0.5)

			r1 := clamp(math.Cos(fx_1*9.0+fy*9.5)*0.5 + 0.5)
			g1 := clamp(math.Cos(fx_1*12.0+fy*40.5)*0.5 + 0.5)
			b1 := clamp(math.Cos(fx_1*23.0+fy*15.5)*0.5 + 0.5)

			src[i+0] = uint8(math.Max(0, math.Min(255, float64(((112*b0-87*g0-26*r0)>>8)+128))))
			src[i+1] = uint8(math.Max(0, math.Min(255, float64(((16*b0+157*g0+47*r0)>>8)+16))))
			src[i+2] = uint8(math.Max(0, math.Min(255, float64(((112*r1-10*b1-102*g1)>>8)+128))))
			src[i+3] = uint8(math.Max(0, math.Min(255, float64(((16*b1+157*g1+47*r1)>>8)+16))))
		}
	}
	return src
}

func transmitFrame(sender *NDI.SendInstanceT, src *byte) {
	// We are going to create a 1920x1080 interlaced frame at 29.97Hz.
	videoFrame := &NDI.VideoFrameV2T{
		Xres:   xRes,
		Yres:   yRes,
		FourCC: NDI.FourCCTypeUYVY,
		PData:  src,
		//LineStrideInBytes: xRes * 2,
		FrameRateN: framerateN,
		FrameRateD: framerateD,
	}

	NDI.SendSendVideoV2(sender, videoFrame)
	videoFrame.Free() // manually free the underlying reference at this point to prevent a memory leak
}

/**
A mildly adopted version of the benchmark application from the standard NewTek examples
*/
func main() {
	if !NDI.Initialize() {
		panic("cannot run NDI")
	}

	sender := NDI.SendCreate(&NDI.SendCreateT{
		PNdiName:   "benchmark",
		ClockAudio: false,
		ClockVideo: false,
	})
	if sender == nil {
		panic("unable to create sender")
	}

	defer func() {
		println("Benchmark stopped")
		NDI.SendSendVideoAsyncV2(sender, nil)
		NDI.SendDestroy(sender)
		NDI.Destroy()
	}()

	// allocate memory
	src := createData()

	prevTime := time.Now()

	println("Running benchmark...")

	// cycle over data
	idx := 0
	for {
		transmitFrame(sender, &src[xRes*2*(idx%(yRes*(scrollDist-1)))])

		if idx != 0 && (idx%1000) == 0 {
			thisTime := time.Now()
			fmt.Printf("%dx%d video encoded at %1.1ffps.\n", xRes, yRes, 1000/thisTime.Sub(prevTime).Seconds())
			prevTime = thisTime
			debug.FreeOSMemory()
		}
		idx++
	}
}
