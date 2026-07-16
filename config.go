package atl

import (
	"fmt"
	"time"
)

var bands = [][4]int{
	{40, 48, 56, 64},
	{80, 88, 96, 104},
	{120, 128, 136, 144},
	{160, 168, 176, 184},
	{200, 208, 216, 224},
	{240, 248, 256, 264},
}

func NewATLConfig(sampleRate int, band int) (ATLConfig, error) {
	if band < 0 || band >= len(bands) {
		return ATLConfig{}, fmt.Errorf("invalid band: %v (must be 0-%v)", band, len(bands)-1)
	}
	return ATLConfig{
		sampleRate,
		bands[band][0], bands[band][1], bands[band][2], bands[band][3],
		1024,
		6, 1,
		10,
	}, nil
}

type ATLConfig struct {
	SampleRate         int
	FrequencyBinAHigh  int
	FrequencyBinALow   int
	FrequencyBinBHigh  int
	FrequencyBinBLow   int
	SamplesPerWindow   int
	WindowsPerBitWrite int
	WindowsPerBitRead  int
	Threshold          float64
}

func (a ATLConfig) CalculateNoteDuration() time.Duration {
	return time.Microsecond * time.Duration(1000000*float64(a.SamplesPerWindow*a.WindowsPerBitWrite)/float64(a.SampleRate))
}

func (a ATLConfig) CalculateSyncDuration() time.Duration {
	return time.Microsecond * time.Duration(1000000*float64(a.SamplesPerWindow*4)/float64(a.SampleRate))
}

func (a ATLConfig) CalculateBitrate() float64 {
	return 1.0 / a.CalculateNoteDuration().Seconds()
}

func (a ATLConfig) CalculateFrequencies() [4]float64 {
	fs := EmptyFreqs(a.SamplesPerWindow, a.SampleRate)
	return [4]float64{
		fs.Freq(int(a.FrequencyBinAHigh)),
		fs.Freq(int(a.FrequencyBinALow)),
		fs.Freq(int(a.FrequencyBinBHigh)),
		fs.Freq(int(a.FrequencyBinBLow)),
	}
}

func (a ATLConfig) CalculateTransmissionDuration(numBytes int) time.Duration {
	return a.CalculateNoteDuration()*time.Duration(numBytes*8) + a.CalculateSyncDuration()
}
