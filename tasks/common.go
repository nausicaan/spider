package tasks

import (
	"fmt"
	"os/exec"
)

var (
	stagingURL, stagingPATH, stagingID string
	prodURL, prodPATH, prodID          string
	siteNAME                           string
)

// Prepare
func Prepare() {
	// exec.Command("cd", "/data/www-app/staging_blog_gov_bc_ca").Run()
	exec.Command("cd", "/Users/byron/Git/go").Run()
	fmt.Print(exec.Command("pwd").Run())
}
