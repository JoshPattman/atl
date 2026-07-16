package atl

import (
	"iter"
	"math/cmplx"

	"github.com/argusdusty/gofft"
)

func AnalyseFreqs(samples []float32, sampleRate int) (*Freqs, error) {
	samplesc128 := make([]complex128, len(samples))
	for i, f := range samples {
		samplesc128[i] = complex(float64(f), 0)
	}
	err := gofft.FFT(samplesc128)
	if err != nil {
		return nil, err
	}
	floatVals := make([]float64, len(samplesc128)/2)
	for i := range floatVals {
		floatVals[i] = cmplx.Abs(samplesc128[i])
	}
	return &Freqs{floatVals, sampleRate}, nil
}

func EmptyFreqs(numSamples int, sampleRate int) *Freqs {
	return &Freqs{make([]float64, numSamples/2), sampleRate}
}

type Freqs struct {
	values     []float64
	sampleRate int
}

func (f *Freqs) All() iter.Seq2[float64, float64] {
	return func(yield func(freq float64, val float64) bool) {
		for i, v := range f.values {
			if !yield(f.Freq(i), v) {
				return
			}
		}
	}
}

func (f *Freqs) Len() int {
	return len(f.values)
}

func (f *Freqs) Value(i int) float64 {
	return f.values[i]
}

func (f *Freqs) Freq(i int) float64 {
	return float64(i*f.sampleRate) / (2 * float64(len(f.values)))
}
