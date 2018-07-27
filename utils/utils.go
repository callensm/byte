package utils

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var logger = NewLogger()
var spinners = make(map[string]*spinner.Spinner)

// Catch handles Go library errors
func Catch(err error) {
	if err != nil {
		logger.Error(err.Error())
	}
}

// CreateSpinner instantiates and starts a new spinner in the terminal
// and adds it to the map of active spinners
func CreateSpinner(set int, clr, text, tag string) {
	if _, ok := spinners[tag]; !ok {
		s := spinner.New(spinner.CharSets[set], 100*time.Millisecond)
		s.Color(clr, "bold")
		s.Suffix = " " + color.WhiteString(text)
		s.Start()
		spinners[tag] = s
	}
}

// RemoveSpinner stops and deletes a spinner from the map
func RemoveSpinner(tag, text string) {
	if s, ok := spinners[tag]; ok {
		s.FinalMSG = fmt.Sprintf("%s %s\n", color.GreenString("✔"), text)
		s.Stop()
		delete(spinners, tag)
	}
}
