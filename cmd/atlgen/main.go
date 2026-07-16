package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/JoshPattman/atl"

	"github.com/jfreymuth/pulse"
	"github.com/spf13/cobra"
)

var atlGen = &cobra.Command{
	Use:   "atlgen",
	Short: "Utilities for generating and parsing ATL (Audio Transport Language).",
}

var (
	atlGenBand int
)

func main() {
	if err := atlGen.Execute(); err != nil {
		fatal(err)
	}
}

// micStream, err := client.NewRecord(
// 	atl.PulseWriter(micBuf),
// 	pulse.RecordSource(source),
// 	pulse.RecordSampleRate(sampleRate),
// 	pulse.RecordLatency(.1),
// 	pulse.RecordMono,
// )

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Fatal: '%s'\n", err.Error())
	os.Exit(1)
}

func catchSignals() <-chan os.Signal {
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	return terminate
}

func getConfig() (atl.ATLConfig, error) {
	return atl.NewATLConfig(48000, atlGenBand)
}

func getClient() (*pulse.Client, error) {
	return pulse.NewClient()
}
