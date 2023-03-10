package tasks

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	automatic string = "\033[0m"
	bgRed     string = "\033[41m"
	fgYellow  string = "\033[33m"
	halt      string = "program halted "
	huh       string = "Unrecognized flag detected -"
)

var reader = bufio.NewReader(os.Stdin)

// Get user input via screen prompt
func solicit(prompt string) string {
	fmt.Print(prompt)
	response, _ := reader.ReadString('\n')
	return strings.TrimSpace(response)
}

// Run standard terminal commands and display the output
func verbose(name string, task ...string) {
	path, err := exec.LookPath(name)
	osCmd := exec.Command(path, task...)
	osCmd.Stdout = os.Stdout
	osCmd.Stderr = os.Stderr
	err = osCmd.Run()
	inspect(err)
}

// Run standard terminal commands and display the output
func silent(name string, task ...string) {
	path, err := exec.LookPath(name)
	inspect(err)
	err = exec.Command(path, task...).Run()
}

// Run a terminal command, then capture and return the output as a byte variable
func byteme(name string, task ...string) []byte {
	path, err := exec.LookPath(name)
	inspect(err)
	osCmd, _ := exec.Command(path, task...).CombinedOutput()
	return osCmd
}

// Read any file and return the contents as a byte variable
func readit(file string) []byte {
	mission, err := os.Open(file)
	inspect(err)
	outcome, err := io.ReadAll(mission)
	inspect(err)
	defer mission.Close()
	return outcome
}

// Check for errors, print the result if found
func inspect(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Provide and highlight an informational message
func banner(message string) {
	fmt.Println(fgYellow)
	fmt.Println("**", automatic, message, fgYellow, "**", automatic)
}

// Tell the program what to do based on the results of a --dry-run
func direct(answer, nav string) {
	if strings.ToLower(answer) == "y" {
		proceed(nav)
	} else {
		os.Exit(0)
	}
}

// Execute the functions without a --dry-run condition
func proceed(action string) {
	switch action {
	case "lf":
		linkFix()
	case "ac":
		assetCopy(sourceOBJ.BlogID, destOBJ.BlogID)
	case "fr":
		uploadsFolder(sourceOBJ.BlogID, destOBJ.BlogID)
	case "fr2":
		uploadsFolderEscapes(sourceOBJ.BlogID, destOBJ.BlogID)
	case "hf":
		httpFind()
	}
}

// Solicite user confirmation after completion of a --dry-run
func confirm(d string) string {
	fmt.Println(d)
	answer := solicit("Does this output seem acceptable, shall we continue without the --dry-run flag? (y/n) ")
	return answer
}

// Alert prints a colourized error message
func Alert(message string) {
	fmt.Println(bgRed, message, halt)
	fmt.Println(automatic)
}
