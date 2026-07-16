package atl

import (
	"io"
	"math"
)

func NewATLWriter(config ATLConfig, sampleWriter SampleWriter) io.Writer {
	return &atlWriter{
		config,
		sampleWriter,
	}
}

type atlWriter struct {
	config ATLConfig
	writer SampleWriter
}

func (a *atlWriter) Write(data []byte) (int, error) {
	samples := a.generate(data)
	_, err := a.writer.Write(samples)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (a *atlWriter) generate(data []byte) []float32 {
	noteDatas := a.toNoteData(data)
	result := make([]float32, len(noteDatas)*a.config.SamplesPerWindow*a.config.WindowsPerBitWrite)
	freqs := a.config.CalculateFrequencies()
	for i := range result {
		noteData := noteDatas[i/(a.config.SamplesPerWindow*a.config.WindowsPerBitWrite)]
		t := float64(i) / float64(a.config.SampleRate)
		result[i] = float32(sample(noteData, freqs, t))
	}
	return result
}

func (a *atlWriter) toNoteData(bs []byte) []noteData {
	notes := make([]noteData, 8*len(bs))
	flipFlop := true
	for offset, b := range bs {
		for i := 7; i >= 0; i-- {
			noteIndex := offset*8 + (7 - i)
			isOn := (b & (1 << i)) != 0
			if flipFlop {
				notes[noteIndex] = noteData{isOn, !isOn, false, false}
			} else {
				notes[noteIndex] = noteData{false, false, isOn, !isOn}
			}
			flipFlop = !flipFlop
		}
	}
	preamble := []noteData{{true, true, true, true}, {true, true, true, true}, {}, {}}
	return append(preamble, notes...)
}

func sample(note noteData, freqs [4]float64, t float64) float64 {
	v := 0.0
	if note.FAHigh {
		v += math.Sin(t * math.Pi * 2 * freqs[0])
	}
	if note.FALow {
		v += math.Sin(t * math.Pi * 2 * freqs[1])
	}
	if note.FBHigh {
		v += math.Sin(t * math.Pi * 2 * freqs[2])
	}
	if note.FBLow {
		v += math.Sin(t * math.Pi * 2 * freqs[3])
	}
	return v / 4
}
