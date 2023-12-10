package main

// Constant declarations
const (
	few  string = "Insufficient arguments supplied -"
	many string = "Too many arguments supplied -"
)

// Start of the Spider application
func main() {
	if len(flag) < 3 {
		alert(few)
	} else if len(flag) > 3 {
		alert(many)
	} else {
		quarterback()
	}
}
