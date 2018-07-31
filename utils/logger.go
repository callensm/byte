package utils

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Logger is an empty struct for use in throughout
// the CLI for colored text outputs
type Logger struct{}

// NewLogger creates and returns a new instance of
// the Logger struct as a pointer
func NewLogger() *Logger {
	return new(Logger)
}

// Info is used for logging basic information to
// the console for visual purposes
func (l *Logger) Info(str string) {
	header := color.BlueString("[INFO]")
	text := color.WhiteString(str)
	display(header, text)
}

// Warn is used for logging warnings to the user
// in the console if something unexpected happened
func (l *Logger) Warn(str string) {
	header := color.YellowString("[WARN]")
	text := color.WhiteString(str)
	display(header, text)
}

// Error is used for logging error messages
// to the console is something went wrong with
// processing or the execution of a command
// based on argument or option inputs
func (l *Logger) Error(str string) {
	header := color.RedString("[ERROR]")
	text := color.WhiteString(str)
	display(header, text)
	os.Exit(1)
}

// FileSuccess is a log that occurs when a file
// has been successfully sent from origin to destination
// and download is complete at its destination
func (l *Logger) FileSuccess(name string) {
	text := color.WhiteString("%s", name)
	display(" ↳ 📄", fmt.Sprintf("%s %s", text, color.GreenString("✔")))
}

// FileError is a log for when a file failed
// to either be sent from origin to destination or the
// file could not successfully download
func (l *Logger) FileError(name string) {
	text := color.WhiteString("%s failed to send", name)
	display(" ↳ 📄", fmt.Sprintf("%s %s", text, color.RedString("𝘅")))
}

// Directory logs which directory the files are coming from and how many
func (l *Logger) Directory(size int, path string, sending bool) {
	var text string
	if sending {
		text = color.WhiteString("Sending %d files from %s:", size, path)
	} else {
		text = color.WhiteString("Writing %d files to %s:", size, path)
	}
	display("📁", text)
}

// Tree is a specific logging header type for displaying JSON encoded
// Tree type structures
func (l *Logger) Tree(tree string) {
	bookend := color.GreenString("[TREE:%d]", len(tree))
	text := color.WhiteString(tree)
	display(bookend, fmt.Sprintf("\n%s\n%s", text, bookend))
}

// Clear prints two special unicode character
// sequences to clear the terminal and move the
// cursor back to the home position
func (l *Logger) Clear() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

// display writes the argued header and text to the console
func display(header, text string) {
	fmt.Printf("%s %s\n", header, text)
}
