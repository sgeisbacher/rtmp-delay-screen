package webrtcutils

import (
	"sync"

	"github.com/pion/webrtc/v3"
)

type TrackProvider struct {
	mx         sync.Mutex
	videoTrack *webrtc.TrackLocalStaticSample
}

func (tp *TrackProvider) Set(videoTrack *webrtc.TrackLocalStaticSample) {
	tp.mx.Lock()
	defer tp.mx.Unlock()

	tp.videoTrack = videoTrack
}

func (tp *TrackProvider) Get() *webrtc.TrackLocalStaticSample {
	tp.mx.Lock()
	defer tp.mx.Unlock()

	return tp.videoTrack
}
