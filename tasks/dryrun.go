package tasks

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// ### Dry Run Functions ###

func linkFixDR() string {
	dr, _ := exec.Command("wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", sourceURL, destURL, "--dry-run").Output()
	return string(dr)
}

func assetCopyDR(sid, did string) string {
	dr, _ := exec.Command("rsync", "-a", "/data/www-assets/"+sourcePath+"/uploads/sites/"+sid+"/", "/data/www-assets/"+destPath+"/uploads/sites/"+did+"/", "--stats", "--dry-run").Output()
	return string(dr)
}

// Correct the uploads folder references
func folderRefDR(sid, did string) string {
	dr, _ := exec.Command("wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", "app/uploads/sites/"+sid, "app/uploads/sites/"+did, "--dry-run").Output()
	return string(dr)
}

func httpFindDR() string {
	dr, _ := exec.Command("wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", "http://", "https://", "--dry-run").Output()
	return string(dr)
}

func direct(answer, nav string) {
	if answer == "Y" {
		proceed(nav)
	}
}

func proceed(action string) {
	switch action {
	case "lf":
		linkFix()
	case "ac":
		assetCopy(sourceOBJ.BlogID, destOBJ.BlogID)
	case "fr":
		folderRef(sourceOBJ.BlogID, destOBJ.BlogID)
	case "fr2":
		folderRef2(sourceOBJ.BlogID, destOBJ.BlogID)
	case "hf":
		httpFind()
	}
}

func confirm(d string) string {
	fmt.Println(d)
	answer := getInput("Does this output seem acceptable, shall we continue without the --dry-run flag?")
	return answer
}

// The getInput function takes a string prompt and asks the user for input.
func getInput(prompt string) string {
	fmt.Print("\n ", prompt)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)
	return userInput
}
