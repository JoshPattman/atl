package main

import (
	"errors"
	"strings"

	"github.com/jfreymuth/pulse"
)

func findSource(client *pulse.Client, search string) (*pulse.Source, error) {
	allSources, err := client.ListSources()
	if err != nil {
		return nil, err
	}
	if len(allSources) == 0 {
		return nil, errors.New("there were no microphones available")
	}
	source := allSources[0]
	if search != "" {
		for _, s := range allSources {
			if strings.Contains(s.Name(), search) {
				source = s
				break
			}
		}
	}
	return source, nil
}
