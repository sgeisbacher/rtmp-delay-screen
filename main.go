package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/pion/webrtc/v3"
	"github.com/sgeisbacher/go-rtmp-screen/ui"
)

func main() {
	http.Handle("/", templ.Handler(ui.PlayerLayout()))
	http.HandleFunc("/createPeerConnection", createPeerConnection)

	fmt.Println("Listening on :8080")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatalf("unknown error: %v\n", err)
	}
}

// Add a single video track
func createPeerConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming HTTP Request")

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

	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypePCMA}, "audio", "pion")
	if err != nil {
		panic(err)
	}
	if _, err = peerConnection.AddTrack(audioTrack); err != nil {
		panic(err)
	}

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

	go startRTMPServer(peerConnection, videoTrack, audioTrack)
}
