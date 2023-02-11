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

/* ##### Dry Run Functions ##### */

// Correct the links with search-replace --dry-run
func linkFixDR() string {
	dr, _ := exec.Command("wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", sourceURL, destURL, "--dry-run").Output()
	return string(dr)
}

// Copy the site assets over --dry-run
func assetCopyDR(sid, did string) string {
	dr, _ := exec.Command("rsync", "-a", "/data/www-assets/"+sourcePath+"/uploads/sites/"+sid+"/", "/data/www-assets/"+destPath+"/uploads/sites/"+did+"/", "--stats", "--dry-run").Output()
	return string(dr)
}

// Correct the references to the uploads folder --dry-run
func uploadsFolderDR(sid, did string) string {
	dr, _ := exec.Command("wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", "app/uploads/sites/"+sid, "app/uploads/sites/"+did, "--dry-run").Output()
	return string(dr)
}

// Correct any unescaped folders due to Gutenberg Blocks --dry-run
func uploadsFolderEscapesDR(sid, did string) string {
	dr, _ := exec.Command("wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", "app\\/uploads\\/sites\\/"+sid, "app\\/uploads\\/sites\\/"+did, "--dry-run").Output()
	return string(dr)
}

// Catch any lingering http addresses --dry-run
func httpFindDR() string {
	dr, _ := exec.Command("wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", "http://", "https://", "--dry-run").Output()
	return string(dr)
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
	answer := getInput("Does this output seem acceptable, shall we continue without the --dry-run flag?")
	return answer
}

// Get user input via screen prompt
func getInput(prompt string) string {
	fmt.Print("\n ", prompt)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)
	return userInput
}
