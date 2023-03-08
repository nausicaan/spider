package tasks

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	// Colour palette.
	colorReset     = "\033[0m"
	fgBrightYellow = "\033[93m"
)

var reader = bufio.NewReader(os.Stdin)

// Get user input via screen prompt
func converse(prompt string) string {
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
	problems(err)
}

// Run a terminal command, then capture and return the output as a byte
func byteme(name string, task ...string) []byte {
	path, err := exec.LookPath(name)
	problems(err)
	osCmd, _ := exec.Command(path, task...).CombinedOutput()
	return osCmd
}

// Check for errors, print the result if found
func problems(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

// The title function displays a section header
func title(content string) {
	fmt.Println(fgBrightYellow, "+-+ +-+ +-+ +-+ +-+ +-+ +-+ +-+ +-+ +-+ +-+")
	fmt.Println(content)
	fmt.Println("+-+ +-+ +-+ +-+ +-+ +-+ +-+ +-+ +-+ +-+ +-+" + colorReset)
}

// Tell the program what to do based on the results of a --dry-run
func direct(answer, nav string) {
	if answer == "Y" {
		proceed(nav)
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
	answer := converse("Does this output seem acceptable, shall we continue without the --dry-run flag?")
	return answer
}
