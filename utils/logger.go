package utils

import (
	"fmt"
	"os"
	"strings"

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

// Clear prints two special unicode character
// sequences to clear the terminal and move the
// cursor back to the home position
func (l *Logger) Clear() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

// Info is used for logging basic information to
// the console for visual purposes
func (l *Logger) Info(str string) {
	h := color.New(color.FgBlue, color.Bold)
	header := h.Sprint("ⓘ ")
	text := color.WhiteString(str)
	display(header, text+"\n")
}

// Warn is used for logging warnings to the user
// in the console if something unexpected happened
func (l *Logger) Warn(str string) {
	h := color.New(color.FgYellow, color.Bold)
	header := h.Sprint("! ")
	text := color.WhiteString(str)
	display(header, text+"\n")
}

// Error is used for logging error messages
// to the console is something went wrong with
// processing or the execution of a command
// based on argument or option inputs
func (l *Logger) Error(str string) {
	h := color.New(color.FgRed, color.Bold)
	header := h.Sprint("𝗫 ")
	text := color.WhiteString(str)
	display(header, text+"\n")
	os.Exit(1)
}

// Prompt requests user import for a given prompt
func (l *Logger) Prompt(str string) string {
	q := color.New(color.FgYellow, color.Bold)
	display(q.Sprint("? "), str)
	res := make([]byte, 2)
	_, err := os.Stdin.Read(res)
	Catch(err)
	return strings.ToLower(strings.Trim(string(res), "\n"))
}

// FileSuccess is a log that occurs when a file
// has been successfully sent from origin to destination
// and download is complete at its destination
func (l *Logger) FileSuccess(name string) {
	text := color.WhiteString("%s", name)
	display(" ↳ 📄", fmt.Sprintf("%s %s\n", text, color.GreenString("✔")))
}

// FileError is a log for when a file failed
// to either be sent from origin to destination or the
// file could not successfully download
func (l *Logger) FileError(name string) {
	text := color.WhiteString("%s failed to send", name)
	display(" ↳ 📄", fmt.Sprintf("%s %s\n", text, color.RedString("𝘅")))
}

// Directory logs which directory the files are coming from and how many
func (l *Logger) Directory(size int, path string, sending bool) {
	var text string
	if sending {
		text = color.WhiteString("Sending %d files from %s:\n", size, path)
	} else {
		text = color.WhiteString("Writing %d files to %s:\n", size, path)
	}
	display("📁", text)
}

// Tree is a specific logging header type for displaying JSON encoded
// Tree type structures
func (l *Logger) Tree(tree string) {
	bookend := color.GreenString("[TREE:%d]", len(tree))
	text := color.WhiteString(tree)
	display(bookend, fmt.Sprintf("\n%s\n%s\n", text, bookend))
}

// display writes the argued header and text to the console
func display(header, text string) {
	fmt.Printf("%s %s", header, text)
}
