package tasks

import (
	"encoding/json"
	"log"
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

var (
	testList, stageList, prodList []Blog
	testID, stageID, prodID       string
	test, stage, prod             Blog

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
		stageList = parseJSON(stageURL, stagePath)
		stage = aquireID("https://"+stageURL+"/"+site+"/", stageList)
		backupDB(stagePath)
	case "p2s":
		prodList = parseJSON(prodURL, prodPath)
		prod = aquireID("https://"+prodURL+"/"+site+"/", prodList)
		backupDB(prodPath)
	default:
		testList = parseJSON(testURL, testPATH)
		test = aquireID("http://test.engage.gov.bc.ca/"+site+"/", testList)
		exportDB(test.URL)
		flush()
	}
}

// Query WordPress for a list of all sites and map the json data to a struct array
func parseJSON(url, path string) []Blog {
	var blog []Blog
	query, _ := exec.Command("wp", "site", "list", "--path=/data/www-app/"+path+"/current/web/wp", "--url="+url, "--format=json").Output()
	json.Unmarshal(query, &blog)
	return blog
}

// Search the blog structure to find the ID that matches the supplied URL
func aquireID(url string, blogs []Blog) Blog {
	var blog Blog
	for _, item := range blogs {
		if item.URL == url {
			blog.BlogID = item.BlogID
			blog.LastUpdated = item.LastUpdated
			blog.Registered = item.Registered
			blog.URL = item.URL
		}
	}
	return blog
}

// Export the database tables
func exportDB(furl string) {
	c1, err := exec.Command("wp db tables", "--url="+furl+"--all-tables-with-prefix --format=csv").Output()
	errors(err)
	exec.Command("wp", "db", "export", "--tables=$("+string(c1)+")", "--quiet", "/data/temp/"+site+".sql").Run()
	// exec.Command("wp", "db", "export", "--tables=$(wp db tables", "--url="+furl, "--all-tables-with-prefix", "--format=csv)", "/data/temp/"+site+".sql").Run()
}

// Create a user export file
func exportUsers(furl, path string) {
	exec.Command("/data/scripts/user_export.py", "-p", "/data/www-app/"+path+"/current/web/wp", "-u", furl, "-o", "/data/temp/"+site+".json").Run()
}

func createSite(url, title, email string) {
	exec.Command("wp", "site", "create", "--url="+url, "--title="+title, "--email="+email, "--quiet").Run()
}

// Import the data
func importDB() {
	exec.Command("wp", "db", "import", "/data/temp/"+site+".sql", "--quiet").Run()
}

// Backup the database
func backupDB(path string) {
	exec.Command("wp", "db", "export", "--path=/data/www-app/"+path+"/current/web/wp", "/data/temp/backup.sql", "--quiet").Run()
}

// Take the blog_id from (fid) the old site and send it to (tid) the new one to be replaced
func replaceIDs(fid, tid string) {
	exec.Command("sed", "-i", "'s/wp_"+fid+"_/wp_"+tid+"_/g'", "/data/temp/"+site+".sql").Run()
}

// Copy the site assets over
func assetCopy(fpath, fid, tpath, tid string) {
	exec.Command("rsync", "-a", "/data/www-assets/"+fpath+"/uploads/sites/"+fid+"/", "/data/www-assets/"+tpath+"/uploads/sites/"+tid+"/").Run()
}

// Correct the links with search-replace
func linkFix(furl, turl string) {
	exec.Command("wp", "search-replace", "--url="+turl, "--all-tables-with-prefix", furl, turl, "--quiet").Run()
}

// Catch any lingering http addresses
func httpFind(turl string) {
	exec.Command("wp", "search-replace", "--url="+turl, "--all-tables-with-prefix", "http://", "https://", "--quiet").Run()
}

// Correct the uploads folder references
func folderRef(turl, fid, tid string) {
	exec.Command("wp", "search-replace", "--url="+turl, "--all-tables-with-prefix", "app/uploads/sites/"+fid, "app/uploads/sites/"+tid, "--quiet").Run()
}

// Remap the users to match their new ID
func remap(turl string) {
	exec.Command("user_import.py", "-p", "current/web/wp", "-u", turl, "-i ", "/data/temp/"+site+".json", "--quiet").Run()
}

// Flush the WordPress cache
func flush() {
	exec.Command("wp", "cache", "flush", "--quiet").Run()
}

// Check for errors, halt the program if found, and log the result
func errors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
