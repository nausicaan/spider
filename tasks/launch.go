package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Blog holds converted json data
type Blog struct {
	BlogID      string `json:"blog_id"`
	URL         string `json:"url"`
	LastUpdated string `json:"last_updated"`
	Registered  string `json:"registered"`
}

const (
	// Test URL's and PATH
	testURL    = "test.blog.gov.bc.ca"
	engTestURL = "test.engage.gov.bc.ca"
	vanTestURL = "test.vanity.blog.gov.bc.ca"
	testPATH   = "test_blog_gov_bc_ca"

	// Production URL's
	stagURL = "staging.blog.gov.bc.ca"
	gwwURL  = "gww.blog.gov.bc.ca"
	prodURL = "engage.gov.bc.ca"

	// Production PATH's
	stagPATH = "staging_blog_gov_bc_ca"
	gwwPATH  = "gww_blog_gov_bc_ca"
	prodPATH = "blog_gov_bc_ca"
)

var (
	testID, stagID, prodID       string
	testBLOG, stagBLOG, prodBLOG []Blog

	flag = os.Args[1]
	site = os.Args[2]
)

/*
Flags:
	s2p - Staging to Production
	p2s - Production to Staging
*/

// Prepare function controls the flow of the program
func Prepare() {
	switch flag {
	case "s2p":
		stagBLOG = parseJSON(stagURL, stagPATH)
		stagID = aquireID("https://"+stagURL+"/"+site+"/", stagBLOG)
	case "p2s":
		prodBLOG = parseJSON(prodURL, prodPATH)
		prodID = aquireID("https://"+prodURL+"/"+site+"/", prodBLOG)
	default:
		testBLOG = parseJSON(testURL, testPATH)
		testID = aquireID("http://test.engage.gov.bc.ca/"+site+"/", testBLOG)
		fmt.Println(testID)
	}
}

// Query WordPress for a list of all sites and map the json data to a struct array
func parseJSON(url, path string) []Blog {
	var blog []Blog
	query, _ := exec.Command("wp", "site", "list", "--path=/data/www-app/"+path+"/current/web/wp", "--url="+url, "--format=json").Output()
	json.Unmarshal(query, &blog)
	// query, _ := exec.Command("wp", "site", "list", "--path=/data/www-app/"+path+"/current/web/wp", "--url="+url).Output()
	// errors(os.WriteFile("site-list.json", query, 0644))
	// grep, _ = exec.Command("grep", site, "site-list.txt").Output()
	// before, _, _ := strings.Cut(string(grep), "h")
	// id := strings.TrimSpace(before)
	return blog
}

// Search the blog structure to find the ID that matches the supplied URL
func aquireID(url string, blog []Blog) string {
	var id string
	for _, item := range blog {
		if item.URL == url {
			id = item.BlogID
		}
	}
	return id
}

// Export the database tables
func exportDB(furl string) {
	exec.Command("wp", "db", "export", "--tables=$(wp db tables", "--url="+furl+"/"+site, "--all-tables-with-prefix", "--format=csv)", "/data/temp/"+site+".sql").Run()
}

// Create a user export file
func exportUsers(furl string) {
	exec.Command("user_export.py", "-p", "current/web/wp", "-u", furl+"/"+site, "-o", "/data/temp/"+site+".json").Run()
}

// Import the data
func importDB() {
	exec.Command("wp", "db", "import", "/data/temp/"+site+".sql").Run()
}

// Backup the database
func backupDB(path string) {
	exec.Command("wp", "db", "export", "--path=/data/www-app/"+path, "/data/temp/backup.sql").Run()
}

// Take the blog_id from (fid) the old site and send it to (tid) the new one to be replaced
func replaceIDs(fid, tid string) {
	exec.Command("sed", "-i", "'s/wp_"+fid+"_/wp_"+tid+"_/g'", "/data/temp/"+site+".sql").Run()
}

// Copy the site assets over
func assetCopy(fpath, fid, tpath, tid string) {
	exec.Command("rsync", "-a", "/data/www-assets/"+fpath+"/uploads/sites/"+fid+"/", "/data/www-assets/"+tpath+"/uploads/sites/"+tid+"/", "--stats").Run()
}

// Correct the links with search-replace
func linkFix(furl, turl string) {
	exec.Command("wp", "search-replace", "--url="+turl+"/"+site, "--all-tables-with-prefix", furl, turl).Run()
}

// Catch any lingering http addresses
func httpFind(turl string) {
	exec.Command("wp", "search-replace", "--url="+turl+"/"+site, "--all-tables-with-prefix", "http://", "https://").Run()
}

// Correct the uploads folder references
func folderRef(turl, fid, tid string) {
	exec.Command("wp", "search-replace", "--url="+turl+"/"+site, "--all-tables-with-prefix", "app/uploads/sites/"+fid, "app/uploads/sites/"+tid).Run()
}

// Remap the users to match their new ID
func remap(turl string) {
	exec.Command("user_import.py", "-p", "current/web/wp", "-u", turl+"/"+site, "-i ", "/data/temp/"+site+".json").Run()
}

// Flush the WordPress cache
func flush() {
	exec.Command("wp", "cache", "flush").Run()
}
