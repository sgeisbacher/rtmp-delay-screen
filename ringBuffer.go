package main

type RingBuffer struct {
	capacity   int
	rbWriteIdx int64
	buffer     [][]byte
}

func createRingBuffer(capacity int) RingBuffer {
	return RingBuffer{
		capacity:   capacity,
		rbWriteIdx: 0,
		buffer:     make([][]byte, capacity),
	}
}

func (r *RingBuffer) Write(data []byte) {
	writeIdx := r.rbWriteIdx % int64(r.capacity)
	r.buffer[writeIdx] = data
	r.rbWriteIdx += 1
}

func (r RingBuffer) Read() ([]byte, bool) {
	if (r.rbWriteIdx + 1) < int64(r.capacity) {
		return nil, false
	}
	readIdx := (r.rbWriteIdx + 1) % int64(r.capacity)
	return r.buffer[readIdx], true
}
