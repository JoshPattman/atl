package atl

import (
	"os"
	"slices"
	"sync"

	"github.com/jfreymuth/pulse"
)

// Return a new buffer that never terminates,
// useful to connect two streams (e.g. a microphone to a reader)
func NewSampleBuffer() *SampleBuffer {
	sb := &SampleBuffer{
		make(chan []float32),
		make(chan readChunkRequest),
		make(chan chan struct{}),
		make(chan struct{}),
		&sync.Once{},
	}
	go sb.run()
	return sb
}

type readChunkRequest struct {
	into   []float32
	result chan int
}

type SampleBuffer struct {
	writeChunk         chan []float32
	readChunk          chan readChunkRequest
	awaitEmptyRequests chan chan struct{}
	closed             chan struct{}
	closedOnce         *sync.Once
}

func (b *SampleBuffer) Write(p []float32) (int, error) {
	select {
	case <-b.closed:
		return 0, os.ErrClosed
	case b.writeChunk <- slices.Clone(p):
		return len(p), nil
	}
}
func (b *SampleBuffer) Read(out []float32) (int, error) {
	select {
	case <-b.closed:
		return 0, pulse.EndOfData
	default:
	}
	request := readChunkRequest{
		out,
		make(chan int),
	}
	select {
	case b.readChunk <- request:
	case <-b.closed:
		return 0, pulse.EndOfData
	}
	select {
	case n := <-request.result:
		return n, nil
	case <-b.closed:
		return 0, pulse.EndOfData
	}
}

func (b *SampleBuffer) run() {
	data := make([]float32, 0)
	for {
		readChunk := b.readChunk
		checkEmpty := b.awaitEmptyRequests
		if len(data) == 0 {
			readChunk = nil
		} else {
			checkEmpty = nil
		}
		select {
		case <-b.closed:
			return
		case chunk := <-b.writeChunk:
			data = append(data, chunk...)
		case request := <-readChunk:
			n := copy(request.into, data)
			data = data[n:]
			select {
			case request.result <- n:
			case <-b.closed:
				return
			}
		case done := <-checkEmpty:
			done <- struct{}{}
		}
	}
}

func (b *SampleBuffer) Close() {
	b.closedOnce.Do(func() {
		close(b.closed)
	})
}

func (b *SampleBuffer) AwaitEmpty() <-chan struct{} {
	s := make(chan struct{}, 1)
	go func() {
		b.awaitEmptyRequests <- s
	}()
	return s
}
