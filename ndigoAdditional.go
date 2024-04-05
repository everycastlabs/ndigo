package ndigo

/*
#cgo CFLAGS: -Wno-deprecated-declarations
#cgo linux LDFLAGS: -L/usr/local/lib -lndi
#cgo darwin LDFLAGS: -Wl,-rpath,/Library/NDI\ SDK\ for\ Apple/lib/macOS -L/Library/NDI\ SDK\ for\ Apple/lib/macOS -lndi
#cgo windows LDFLAGS: -LC:/Program\ Files/NDI/NDI\ 5\ Runtime/v5 -lProcessing.NDI.Lib.x64
#include <stdlib.h>
#include "include/Processing.NDI.Lib.h"
#include "cgo_helpers.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"strconv"
	"unsafe"
)

type SourceType2 struct {
	Name       string
	URLAddress string
}

func FindGetCurrentSources2(instance FindInstanceType) []*SourceType2 {
	var pNoSources C.uint32_t
	pSources := C.NDIlib_find_get_current_sources(C.NDIlib_find_instance_t(instance), &pNoSources)
	if pNoSources == 0 {
		return nil
	}
	sources := (*[1 << 28]C.NDIlib_source_t)(unsafe.Pointer(pSources))[:pNoSources:pNoSources]
	result := make([]*SourceType2, pNoSources)
	for i, source := range sources {
		result[i] = &SourceType2{
			Name:       C.GoString(source.p_ndi_name),
			URLAddress: C.GoString(*(**C.char)(unsafe.Pointer(&source.anon0))),
		}
	}
	return result
}

func FindWaitForSources2(instance FindInstanceType, timeoutInMS uint32) bool {
	return (bool)(C.NDIlib_find_wait_for_sources(C.NDIlib_find_instance_t(instance), C.uint32_t(timeoutInMS)))
}

func SendGetTally2(pInstance SendInstanceType, pTally *TallyType, timeoutInMs uint32) bool {
	var tally C.NDIlib_tally_t
	__ret := C.NDIlib_send_get_tally(C.NDIlib_send_instance_t(pInstance), &tally, 1000)
	//now apply OnPreview and onProgram....
	pTally.OnPreview = (bool)(tally.on_preview)
	pTally.OnProgram = (bool)(tally.on_program)
	__v := (bool)(__ret)
	return __v
}

func RecvGetQueue2(pInstance RecvInstanceType, pTotal *RecvQueueType) {
	var total C.NDIlib_recv_queue_t
	cpInstance, cpInstanceAllocMap := *(*C.NDIlib_recv_instance_t)(unsafe.Pointer(&pInstance)), cgoAllocsUnknown
	C.NDIlib_recv_get_queue(cpInstance, &total)
	runtime.KeepAlive(cpInstanceAllocMap)
	pTotal.VideoFrames = int32(total.video_frames)
	pTotal.AudioFrames = int32(total.audio_frames)
	pTotal.MetadataFrames = int32(total.metadata_frames)
}

type VideoFrameV2Type2 C.NDIlib_video_frame_v2_t

func (t VideoFrameV2Type2) FrameFormatType() int {
	return int(t.frame_format_type)
}

func (t VideoFrameV2Type2) Timecode() int64 {
	return int64(t.timecode * 100)
}

func (t VideoFrameV2Type2) Timestamp() int64 {
	return int64(t.timestamp * 100)
}

func (t VideoFrameV2Type2) FrameFourCC() FourCCVideoType {
	return FourCCVideoType(t.FourCC)
}

func (t VideoFrameV2Type2) Xres() int {
	return int(t.xres)
}

func (t VideoFrameV2Type2) Yres() int {
	return int(t.yres)
}

func (t VideoFrameV2Type2) LineStrideInBytesOrDataSizeInBytes() int {

	//if t.FourCC is not compressed then its line stride...
	//if t.FourCC is compressed then its Data Size
	s := *(**C.int)(unsafe.Pointer(&t.anon0))
	intVar, _ := strconv.Atoi(fmt.Sprintf("%d", s))
	return intVar
}

func (t VideoFrameV2Type2) Data() []byte {
	// lineStrideOrDataSize := t.LineStrideInBytesOrDataSizeInBytes()
	size := (t.Xres() * t.Yres() * 2)
	v := make([]byte, size)

	d := (*[1 << 30]byte)(unsafe.Pointer(t.p_data)) // Read
	copy(v, d[:])
	return v
	// b := v[:lineStrideOrDataSize]
	// return b
	// size := (t.Xres() * t.Yres() * 2)

	//I need to stop C.GoString from doing bad things and rmeoving chars

	//str := C.GoString((*C.char)(unsafe.Pointer(t.p_data)))
	//if the image is pure black.... C.GoString has rather rudely removed those blacks.

	//tmp := []byte(str)
	//return tmp
}

func (t VideoFrameV2Type2) FrameRateN() int {
	return int(t.frame_rate_N)
}

func (t VideoFrameV2Type2) FrameRateD() int {
	return int(t.frame_rate_D)
}

type AudioFrameV3Type2 C.NDIlib_audio_frame_v3_t

func (t AudioFrameV3Type2) NoSamples() int {
	return int(t.no_samples)
}

func (t AudioFrameV3Type2) NoChannels() int {
	return int(t.no_channels)
}

func (t AudioFrameV3Type2) SampleRate() int {
	return int(t.sample_rate)
}

func (t AudioFrameV3Type2) FrameFourCC() FourCCAudioType {
	return FourCCAudioType(t.FourCC)
}

func (t AudioFrameV3Type2) ChannelStrideInBytesOrDataSizeInBytes() int {

	//if t.FourCC is not compressed then its channel stride...
	//if t.FourCC is compressed then its Data Size
	s := *(**C.int)(unsafe.Pointer(&t.anon0))
	intVar, _ := strconv.Atoi(fmt.Sprintf("%d", s))
	return intVar
}

func (t AudioFrameV3Type2) Data() []byte {
	//return []byte(C.GoString((*C.char)(unsafe.Pointer(t.p_data))))
	// samples := t.NoSamples()
	// channels := t.NoChannels()
	channelStride := *(*C.int)(unsafe.Pointer(&t.anon0[0]))
	size := (channelStride * t.no_channels)
	// log.Printf("channelStride: %d, samples: %d, channels: %d, size: %d", channelStride, samples, channels, size)

	return C.GoBytes(unsafe.Pointer(t.p_data), C.int(size))
	// v := make([]byte, size)

	// d := (*[1 << 30]byte)(unsafe.Pointer(t.p_data)) // Read
	// copy(v, d[:])
	// return v
}

func (t AudioFrameV3Type2) Timestamp() int64 {
	return int64(t.timestamp * 100)
}

type MetadataFrameType2 C.NDIlib_metadata_frame_t

func RecvCaptureV32(pInstance RecvInstanceType, pVideoData *VideoFrameV2Type2, pAudioData *AudioFrameV3Type2, pMetadata *MetadataFrameType2, timeoutInMs uint32) FrameType {
	cpInstance, cpInstanceAllocMap := *(*C.NDIlib_recv_instance_t)(unsafe.Pointer(&pInstance)), cgoAllocsUnknown
	// cpVideoData, cpVideoDataAllocMap := pVideoData.PassRef()
	// cpAudioData, cpAudioDataAllocMap := pAudioData.PassRef()
	// cpMetadata, cpMetadataAllocMap := pMetadata.PassRef()
	var videoFrame C.NDIlib_video_frame_v2_t
	var audioFrame C.NDIlib_audio_frame_v3_t
	var metadataFrame C.NDIlib_metadata_frame_t
	ctimeoutInMs, ctimeoutInMsAllocMap := (C.uint32_t)(timeoutInMs), cgoAllocsUnknown
	__ret := C.NDIlib_recv_capture_v3(cpInstance, &videoFrame, &audioFrame, &metadataFrame, ctimeoutInMs)
	runtime.KeepAlive(ctimeoutInMsAllocMap)
	runtime.KeepAlive(cpInstanceAllocMap)
	__v := (FrameType)(__ret)

	if pVideoData != nil {
		*pVideoData = VideoFrameV2Type2(videoFrame)
	}
	if pAudioData != nil {
		*pAudioData = AudioFrameV3Type2(audioFrame)
	}
	if pMetadata != nil {
		*pMetadata = MetadataFrameType2(metadataFrame)
	}

	return __v
}

func RecvFreeVideoV22(instance RecvInstanceType, videoFrame *VideoFrameV2Type2) {
	C.NDIlib_recv_free_video_v2(C.NDIlib_recv_instance_t(instance), (*C.NDIlib_video_frame_v2_t)(videoFrame))
}

func RecvFreeAudioV32(instance RecvInstanceType, audioFrame *AudioFrameV3Type2) {
	C.NDIlib_recv_free_audio_v3(C.NDIlib_recv_instance_t(instance), (*C.NDIlib_audio_frame_v3_t)(audioFrame))
}

func RecvFreeMetadata2(instance RecvInstanceType, metadataFrame *MetadataFrameType2) {
	C.NDIlib_recv_free_metadata(C.NDIlib_recv_instance_t(instance), (*C.NDIlib_metadata_frame_t)(metadataFrame))
}

func GetRedistUrl() string {
	return string(C.NDILIB_REDIST_URL)
}
