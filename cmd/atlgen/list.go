package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var atlGenList = &cobra.Command{
	Use:   "list",
	Short: "List all audio devices that can be used for input and output.",
	Run:   execList,
}

func init() {
	atlGen.AddCommand(atlGenList)
}

func execList(cmd *cobra.Command, args []string) {
	err := runList()
	if err != nil {
		fatal(err)
	}
}

func runList() error {
	client, err := getClient()
	if err != nil {
		return err
	}

	fmt.Println("sources:")
	allSources, err := client.ListSources()
	if err != nil {
		return err
	}
	defaultSource, err := client.DefaultSource()
	if err != nil {
		return err
	}
	for _, s := range allSources {
		fmt.Printf("- id: %v\n  name: %v\n", s.ID(), s.Name())
		if s.ID() == defaultSource.ID() {
			fmt.Println("  default: true")
		}
	}

	fmt.Println("sinks:")
	allSinks, err := client.ListSinks()
	if err != nil {
		return err
	}
	defaultSink, err := client.DefaultSink()
	if err != nil {
		return err
	}
	for _, s := range allSinks {
		fmt.Printf("- id: %v\n  name: %v\n", s.ID(), s.Name())
		if s.ID() == defaultSink.ID() {
			fmt.Println("  default: true")
		}
	}
	return nil
}
