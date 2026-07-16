package atl

import (
	"math"
	"testing"
	"time"

	"github.com/jfreymuth/pulse"
)

func TestSampleBufWriteRead(t *testing.T) {
	t.Log("Init")
	b := NewSampleBuffer()
	n, err := b.Write([]float32{0, 1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	if n != 4 {
		t.Fatal("expected n 4 got", n)
	}
	t.Log("Written")

	a := make([]float32, 3)
	n, err = b.Read(a)

	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatal("expected n 3 got", n)
	}
	if a[0] != 0 || a[1] != 1 || a[2] != 2 {
		t.Fatal("expected 0,1,2 but got", a)
	}
	t.Log("Read 1")

	n, err = b.Read(a)

	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal("expected n 1 got", n)
	}
	if a[0] != 3 {
		t.Fatal("expected 3 but got", a)
	}
	t.Log("Read 2")

	go func() {
		time.Sleep(time.Second)
		b.Close()
	}()
	tStart := time.Now()
	n, err = b.Read(a)
	since := time.Since(tStart)
	if err != pulse.EndOfData {
		t.Fatal("expected end of data err but got", err)
	}
	if n != 0 {
		t.Fatal("expected n 0 got", n)
	}
	if math.Abs((since - time.Second).Seconds()) > 0.1 {
		t.Fatal("expected time 1 got", since)
	}
	t.Log("closed")
	b.Close()
	t.Log("closed 2")
}
