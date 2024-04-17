package main

import (
	"encoding/json"
	"net/http"

	"github.com/pion/webrtc/v3"
	webrtcutils "github.com/sgeisbacher/go-rtmp-screen/webrtc-utils"
)

// Add a single video track
func buildCreatePeerConnectionHandleFunc(videoTrackProvider *webrtcutils.TrackProvider) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
		if err != nil {
			panic(err)
		}

		videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion")
		if err != nil {
			panic(err)
		}
		if _, err = peerConnection.AddTrack(videoTrack); err != nil {
			panic(err)
		}
		videoTrackProvider.Set(videoTrack)

		var offer webrtc.SessionDescription
		if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
			panic(err)
		}

		if err := peerConnection.SetRemoteDescription(offer); err != nil {
			panic(err)
		}

		gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			panic(err)
		} else if err = peerConnection.SetLocalDescription(answer); err != nil {
			panic(err)
		}
		<-gatherComplete

		response, err := json.Marshal(peerConnection.LocalDescription())
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(response); err != nil {
			panic(err)
		}

	}
}
