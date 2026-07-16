package atl

import "github.com/jfreymuth/pulse"

type SampleWriter interface {
	Write(p []float32) (int, error)
}

type SampleReader interface {
	Read(out []float32) (int, error)
}

func PulseWriter(writer SampleWriter) pulse.Writer {
	return pulse.Float32Writer(writer.Write)
}

func PulseReader(reader SampleReader) pulse.Reader {
	return pulse.Float32Reader(reader.Read)
}
