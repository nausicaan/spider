package main

import (
	"os"

	t "github.com/nausicaan/spider/tasks"
)

// Constant declarations
const (
	zero string = "Insufficient arguments supplied -"
	many string = "Too many arguments supplied -"
)

// Start of the Spider application
func main() {
	if len(os.Args) < 3 {
		t.Alert(zero)
	} else if len(os.Args) > 3 {
		t.Alert(many)
	} else {
		t.Quarterback()
	}
}
