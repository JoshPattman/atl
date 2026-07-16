package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/JoshPattman/atl"

	"github.com/jfreymuth/pulse"
	"github.com/spf13/cobra"
)

var atlGenRecv = &cobra.Command{
	Use:   "recv",
	Short: "Receive data from an audio input device.",
	Run:   execRecv,
}
var (
	atlGenRecvSourceID string
)

func init() {
	atlGenRecv.Flags().StringVarP(&atlGenRecvSourceID, "source-id", "i", "", "When specified, will use that source instead of the default.")
	atlGenRecv.Flags().IntVarP(&atlGenBand, "band", "b", 0, "Frequency band to use")
	atlGen.AddCommand(atlGenRecv)
}

func execRecv(cmd *cobra.Command, args []string) {
	err := runRecv(atlGenRecvSourceID)
	if err != nil {
		fatal(err)
	}
}

func runRecv(sourceID string) error {
	terminate := catchSignals()
	output := getOutput()
	config, err := getConfig()
	if err != nil {
		return err
	}
	buf := atl.NewSampleBuffer()
	reader := atl.NewATLReader(config, buf)
	client, err := getClient()
	if err != nil {
		return err
	}
	defer client.Close()
	playbackStream, err := getMicStream(client, buf, config.SampleRate, sourceID)
	if err != nil {
		return err
	}
	defer playbackStream.Close()
	return recvOutputFromAudio(buf, reader, playbackStream, output, terminate)
}

func getOutput() io.Writer {
	return os.Stdout
}

func getMicStream(client *pulse.Client, targetBuf atl.SampleWriter, sampleRate int, sourceID string) (*pulse.RecordStream, error) {
	opts := []pulse.RecordOption{
		pulse.RecordSampleRate(sampleRate),
		pulse.RecordLatency(.1),
		pulse.RecordMono,
	}
	if sourceID != "" {
		sink, err := client.SourceByID(sourceID)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("could not get source with ID '%s'", sourceID), err)
		}
		opts = append(opts, pulse.RecordSource(sink))
	}

	return client.NewRecord(
		atl.PulseWriter(targetBuf),
		opts...,
	)
}

func recvOutputFromAudio(atlBuffer *atl.SampleBuffer, atlReader io.Reader, recordStream *pulse.RecordStream, output io.Writer, terminate <-chan os.Signal) error {
	go recordStream.Start()
	incoming := streamChunks(atlReader)

	finish := func() {
		atlBuffer.Close()
	}

	for {
		select {
		case <-terminate:
			finish()
			return nil
		case chunk := <-incoming:
			_, err := output.Write(chunk)
			if err != nil {
				finish()
				return err
			}
		}
	}
}
