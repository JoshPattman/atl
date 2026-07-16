package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/JoshPattman/atl"

	"github.com/jfreymuth/pulse"
	"github.com/spf13/cobra"
)

var atlGenSend = &cobra.Command{
	Use:   "send",
	Short: "Send data over an audio output device.",
	Run:   execSend,
}
var (
	atlGenSendWithInput string
	atlGenSendSinkID    string
)

func init() {
	atlGenSend.Flags().StringVarP(&atlGenSendWithInput, "message", "m", "", "The message to send. If not specified, will use stdin instead.")
	atlGenSend.Flags().StringVarP(&atlGenSendSinkID, "sink-id", "o", "", "When specified, will use that sink instead of the default.")
	atlGenSend.Flags().IntVarP(&atlGenBand, "band", "b", 0, "Frequency band to use")
	atlGen.AddCommand(atlGenSend)
}

func execSend(cmd *cobra.Command, args []string) {
	err := runSend(atlGenSendWithInput, atlGenSendSinkID)
	if err != nil {
		fatal(err)
	}
}

// The actual run function for send
func runSend(messageToSend string, sinkID string) error {
	terminate := catchSignals()
	messageChunks := getInput(messageToSend)
	config, err := getConfig()
	if err != nil {
		return err
	}
	buf := atl.NewSampleBuffer()
	writer := atl.NewATLWriter(config, buf)
	client, err := getClient()
	if err != nil {
		return err
	}
	defer client.Close()
	speakerStream, err := getSpeakerStream(client, buf, config.SampleRate, sinkID)
	if err != nil {
		return err
	}
	defer speakerStream.Close()
	return playAudioFromInput(buf, writer, speakerStream, messageChunks, terminate)
}

// Get the input source - if a message is specified use that, otherwise use stdin.
func getInput(messageToSend string) <-chan []byte {
	var dataSource io.Reader
	if messageToSend != "" {
		dataSource = bytes.NewReader([]byte(messageToSend))
	} else {
		dataSource = os.Stdin
	}
	return streamChunks(dataSource)
}

// Stream chunks of data out of the blocking reader and into the channel.
func streamChunks(r io.Reader) <-chan []byte {
	messageChunks := make(chan []byte)
	go func() {
		for {
			chunk := make([]byte, 1024)
			n, err := r.Read(chunk)
			if errors.Is(err, io.EOF) {
				close(messageChunks)
				return
			} else if err != nil {
				fmt.Println("Error reading stdin - should not have happened:", err.Error())
				return
			} else if n > 0 {
				messageChunks <- chunk[:n]
			}
		}
	}()
	return messageChunks
}

// Create the speaker stream.
func getSpeakerStream(client *pulse.Client, sourceBuf atl.SampleReader, sampleRate int, sinkID string) (*pulse.PlaybackStream, error) {
	opts := []pulse.PlaybackOption{
		pulse.PlaybackSampleRate(sampleRate),
		pulse.PlaybackLatency(.1),
		pulse.PlaybackMono,
	}
	if sinkID != "" {
		sink, err := client.SinkByID(sinkID)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("could not get sink with ID '%s'", sinkID), err)
		}
		opts = append(opts, pulse.PlaybackSink(sink))
	}

	return client.NewPlayback(
		atl.PulseReader(sourceBuf),
		opts...,
	)
}

// Colpy incoming chunks to the writer, waiting for playback to complete before exiting.
func playAudioFromInput(atlBuffer *atl.SampleBuffer, atlWriter io.Writer, speakerStream *pulse.PlaybackStream, messageChunks <-chan []byte, terminate <-chan os.Signal) error {
	go speakerStream.Start()
	finish := func(now bool) {
		if !now {
			select {
			case <-atlBuffer.AwaitEmpty():
			case <-terminate:
			}
		}
		atlBuffer.Close()
		speakerStream.Drain()
		for speakerStream.Running() {
			time.Sleep(time.Millisecond * 10)
		}
	}

	for {
		select {
		case <-terminate:
			finish(true)
			return nil
		case chunk, ok := <-messageChunks:
			if !ok {
				finish(false)
				return nil
			} else {
				_, err := atlWriter.Write(chunk)
				if err != nil {
					finish(true)
					return err
				}
			}
		}
	}
}
