//go:build ignore

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

type VideoFrameV2Type3 C.NDIlib_video_frame_v2_t

type AudioFrameV3Type3 C.NDIlib_audio_frame_v3_t
