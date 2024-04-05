package ringBuffer

import (
	"sync"
	"time"
)

const MEASURE_TIME_WINDOW = 5

type RingBuffer struct {
	mu          sync.Mutex
	capacity    int
	rbWriteIdx  int64
	buffer      [][]byte
	byteCtr     [][]int
	lastMeasure time.Time
}

func CreateRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		capacity:   capacity,
		rbWriteIdx: 0,
		buffer:     make([][]byte, capacity),
		byteCtr:    make([][]int, MEASURE_TIME_WINDOW),
	}
}

func (r *RingBuffer) GetCapacity() int {
	return r.capacity
}

func (r *RingBuffer) Write(data []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// fmt.Printf("bufsize: %20d (%20d)\n", r.capacity, r.rbWriteIdx)
	measure(r, data)
	writeIdx := r.rbWriteIdx % int64(r.capacity)
	r.buffer[writeIdx] = data
	r.rbWriteIdx += 1
}

func (r *RingBuffer) Read() ([]byte, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if (r.rbWriteIdx + 1) < int64(r.capacity) {
		return nil, false
	}
	readIdx := (r.rbWriteIdx + 1) % int64(r.capacity)
	return r.buffer[readIdx], true
}

func (r *RingBuffer) Reset(newCapacity int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if newCapacity < 0 {
		// keep capacity
		newCapacity = r.capacity
	}
	r.byteCtr = make([][]int, MEASURE_TIME_WINDOW)
	r.lastMeasure = time.Time{}
	r.buffer = make([][]byte, newCapacity)
	r.rbWriteIdx = 0
	r.capacity = newCapacity
}

func (r *RingBuffer) Stats() (dataRate float32, frameRate float32) {
	frameCount := 0
	dataSum := 0
	for _, frames := range r.byteCtr {
		for _, frameBytes := range frames {
			dataSum += frameBytes
			frameCount += 1
		}
	}
	dataRate = float32(dataSum) / MEASURE_TIME_WINDOW
	frameRate = float32(frameCount) / MEASURE_TIME_WINDOW
	return
}

func measure(rb *RingBuffer, data []byte) {
	now := time.Now()
	prevTimeWindow := rb.lastMeasure.Unix() % MEASURE_TIME_WINDOW
	timeWindow := now.Unix() % MEASURE_TIME_WINDOW
	if prevTimeWindow != timeWindow {
		rb.byteCtr[timeWindow] = []int{}
	}
	rb.byteCtr[timeWindow] = append(rb.byteCtr[timeWindow], len(data))
	rb.lastMeasure = now
}
