package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/bep/debounce"
	"github.com/sgeisbacher/go-rtmp-screen/ringBuffer"
	"github.com/sgeisbacher/go-rtmp-screen/ui"
)

const MAX_BUF_SECS = 30
const FRAME_RATE = 30

func main() {
	bufferCapDecouncer := debounce.New(1 * time.Second)
	desiredCapacity := 5 * FRAME_RATE // 5 seconds
	buffer := ringBuffer.CreateRingBuffer(desiredCapacity)

	http.Handle("/", templ.Handler(ui.PlayerLayout()))
	http.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) {
		ui.AdminHomePage(desiredCapacity/FRAME_RATE, MAX_BUF_SECS).Render(r.Context(), w)
	})
	http.HandleFunc("POST /admin/rb/inc/{value}", func(w http.ResponseWriter, r *http.Request) {
		increaseValue, err := strconv.Atoi(r.PathValue("value"))
		if err != nil {
			fmt.Printf("E: invalid increase value: %s", r.PathValue("value"))
			w.WriteHeader(500)
			return
		}
		desiredCapacity += increaseValue * FRAME_RATE
		bufferCapDecouncer(func() { buffer.Reset(desiredCapacity) })

		ui.RingBufferInfos(toSecs(desiredCapacity), toSecs(buffer.GetCapacity()), MAX_BUF_SECS).Render(r.Context(), w)
	})
	http.HandleFunc("GET /streamer/status", func(w http.ResponseWriter, r *http.Request) {
		statusMsg := ""
		switch buffer.Status() {
		case "streaming":
			statusMsg = ""
			break
		case "buffering":
			framesLeft, _ := buffer.BufferingFramesLeft()
			secsLeft := framesLeft / FRAME_RATE
			statusMsg = fmt.Sprintf("%s (%ds) ...", buffer.Status(), secsLeft)
			break
		case "disconnected":
			statusMsg = "disconnected!<br>please (re)start streaming app on phone!"
		default:
			statusMsg = buffer.Status()
		}

		io.WriteString(w, statusMsg)
	})
	http.HandleFunc("GET /admin/infobox/buffer", func(w http.ResponseWriter, r *http.Request) {
		desiredSecs := toSecs(desiredCapacity)
		actualSecs := toSecs(buffer.GetCapacity())
		ui.RingBufferInfos(desiredSecs, actualSecs, MAX_BUF_SECS).Render(r.Context(), w)
	})
	http.HandleFunc("GET /admin/infobox/status", func(w http.ResponseWriter, r *http.Request) {
		status := buffer.Status()
		ui.StatusInfos(status).Render(r.Context(), w)
	})
	http.HandleFunc("GET /admin/infobox/framerate", func(w http.ResponseWriter, r *http.Request) {
		_, frameRate := buffer.Stats()
		ui.FrameRateInfos(frameRate).Render(r.Context(), w)
	})
	http.HandleFunc("GET /admin/infobox/datarate", func(w http.ResponseWriter, r *http.Request) {
		dataRate, _ := buffer.Stats()
		ui.DataRateInfos(dataRate/1024).Render(r.Context(), w)
	})
	http.HandleFunc("/createPeerConnection", buildCreatePeerConnectionHandleFunc(buffer))

	fmt.Println("Listening on :8080")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatalf("unknown error: %v\n", err)
	}
}

func toSecs(n int) int {
	return n / FRAME_RATE
}
