package atl

import (
	"errors"
	"io"
	"time"

	"github.com/jfreymuth/pulse"
)

func NewATLReader(config ATLConfig, sampleReader SampleReader) io.Reader {
	return &atlReader{
		config,
		sampleReader,
		atlReaderState{expectA: true},
	}
}

type atlReader struct {
	config ATLConfig
	reader SampleReader
	state  atlReaderState
}

type atlReaderState struct {
	expectA     bool
	nAHigh      int
	nALow       int
	nBHigh      int
	nBLow       int
	currentBits []bool
}

func (a *atlReader) Read(p []byte) (n int, err error) {
	for i := range p {
		for {
			err := a.readWindowAndUpdateState()
			if err != nil {
				return 0, err
			}
			// fmt.Println(a.state)
			a.updateBits()
			b, ok := a.emitByte()
			if ok {
				p[i] = b
				break
			}
		}
	}
	return len(p), nil
}

func (a *atlReader) updateBits() {
	aHigh := a.state.nAHigh >= a.config.WindowsPerBitRead
	aLow := a.state.nALow >= a.config.WindowsPerBitRead
	bHigh := a.state.nBHigh >= a.config.WindowsPerBitRead
	bLow := a.state.nBLow >= a.config.WindowsPerBitRead

	if aHigh && aLow && bHigh && bLow {
		a.state = atlReaderState{expectA: true}
		// fmt.Println("Read reset")
	} else if a.state.expectA && aHigh && !aLow && !bHigh && !bLow {
		a.state.currentBits = append(a.state.currentBits, true)
		a.state.expectA = !a.state.expectA
		// fmt.Println("Read A high")
	} else if a.state.expectA && !aHigh && aLow && !bHigh && !bLow {
		a.state.currentBits = append(a.state.currentBits, false)
		a.state.expectA = !a.state.expectA
		// fmt.Println("Read A low")
	} else if !a.state.expectA && !aHigh && !aLow && bHigh && !bLow {
		a.state.currentBits = append(a.state.currentBits, true)
		// fmt.Println("Read B high")
		a.state.expectA = !a.state.expectA
	} else if !a.state.expectA && !aHigh && !aLow && !bHigh && bLow {
		a.state.currentBits = append(a.state.currentBits, false)
		a.state.expectA = !a.state.expectA
		// fmt.Println("Read B low")
	}
}

func (a *atlReader) emitByte() (byte, bool) {
	if len(a.state.currentBits) < 8 {
		return 0, false
	}

	var byt byte
	for i, b := range a.state.currentBits[:8] {
		if b {
			byt |= 1 << (7 - i)
		}
	}
	a.state.currentBits = a.state.currentBits[8:]
	return byt, true
}

func (a *atlReader) readWindowAndUpdateState() error {
	fs, err := a.readWindow()
	if err != nil {
		return err
	}
	if fs.Value(a.config.FrequencyBinAHigh) > a.config.Threshold {
		a.state.nAHigh += 1
	} else {
		a.state.nAHigh = 0
	}
	if fs.Value(a.config.FrequencyBinALow) > a.config.Threshold {
		a.state.nALow += 1
	} else {
		a.state.nALow = 0
	}
	if fs.Value(a.config.FrequencyBinBHigh) > a.config.Threshold {
		a.state.nBHigh += 1
	} else {
		a.state.nBHigh = 0
	}
	if fs.Value(a.config.FrequencyBinBLow) > a.config.Threshold {
		a.state.nBLow += 1
	} else {
		a.state.nBLow = 0
	}
	// fmt.Println(a.state)
	return nil
}

func (a *atlReader) readWindow() (*Freqs, error) {
	totalN := 0
	samples := make([]float32, a.config.SamplesPerWindow)
	for totalN != len(samples) {
		n, err := a.reader.Read(samples[totalN:])
		if errors.Is(err, pulse.EndOfData) {
			return nil, io.EOF
		} else if err != nil {
			return nil, err
		}
		totalN += n
		if n == 0 {
			time.Sleep(time.Millisecond)
		}
	}
	return AnalyseFreqs(samples, a.config.SampleRate)
}
