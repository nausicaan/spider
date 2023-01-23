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
	flag = os.Args[1]
	site = os.Args[2]

	testID, stageID, prodID                  string
	testObj, stageObj, prodObj               Blog
	testList, stageList, prodList            []Blog
	sourcePath, sourceURL, destPath, destURL string
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
		sourcePath, sourceURL = stagePath, stageURL
		stageList = parseJSON(sourceURL, sourcePath)
		stageObj = aquireID("https://"+sourceURL+"/"+site+"/", stageList)
		first(stageObj)
		destPath, destURL = prodPath, prodURL
		prodList = parseJSON(destURL, destPath)
		prodObj = aquireID("https://"+destURL+"/"+site+"/", prodList)
		second(stageObj, prodObj, destPath)
	case "p2s":
		sourcePath, sourceURL = prodPath, prodURL
		prodList = parseJSON(sourceURL, sourcePath)
		prodObj = aquireID("https://"+sourceURL+"/"+site+"/", prodList)
		first(prodObj)
		destPath, destURL = stagePath, stageURL
		stageList = parseJSON(destURL, destPath)
		stageObj = aquireID("https://"+destURL+"/"+site+"/", stageList)
		second(prodObj, stageObj, destPath)
	default:
		testList = parseJSON(testURL, testPath)
		testObj = aquireID("http://test.engage.gov.bc.ca/"+site+"/", testList)
		first(testObj)
		// second(testObj, testPath)
	}
}

// Run the first few functions up to the new site creation
func first(source Blog) {
	exportDB(source.URL)
	exportUsers(source.URL, sourcePath)
	createSite("https://"+destURL+"/"+site+"/", site, adminEmail)
}

// Run the remaining functions after being able to grab the new site ID
func second(source, dest Blog, path string) {
	backupDB(path)
	replaceIDs(source.BlogID, dest.BlogID)
	importDB()
	linkFix(source.BlogID, dest.BlogID)
	assetCopy(sourcePath, source.BlogID, destPath, dest.BlogID)
	folderRef(destURL, source.BlogID, dest.BlogID)
	httpFind(destURL)
	remap(destURL)
	flush()
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
func exportDB(surl string) {
	sub, err := exec.Command("wp db tables", "--url="+surl+"--all-tables-with-prefix --format=csv").Output()
	errors(err)
	exec.Command("wp", "db", "export", "--tables=$("+string(sub)+")", "--quiet", "/data/temp/"+site+".sql").Run()
	// exec.Command("wp", "db", "export", "--tables=$(wp db tables", "--url="+surl, "--all-tables-with-prefix", "--format=csv)", "/data/temp/"+site+".sql").Run()
}

// Create a user export file
func exportUsers(surl, path string) {
	exec.Command("/bin/bash", "-c", "/data/scripts/user_export.py", "-p", "/data/www-app/"+path+"/current/web/wp", "-u", surl, "-o", "/data/temp/"+site+".json").Run()
}

func createSite(durl, title, email string) {
	exec.Command("wp", "site", "create", "--url="+durl, "--title="+title, "--email="+email, "--quiet").Run()
}

// Backup the database
func backupDB(path string) {
	exec.Command("wp", "db", "export", "--path=/data/www-app/"+path+"/current/web/wp", "/data/temp/backup.sql", "--quiet").Run()
}

// Take the blog_id from the source (sid) and send it to the destination (did) to be replaced
func replaceIDs(sid, did string) {
	exec.Command("sed", "-i", "'s/wp_"+sid+"_/wp_"+did+"_/g'", "/data/temp/"+site+".sql").Run()
}

// Import the data
func importDB() {
	exec.Command("wp", "db", "import", "/data/temp/"+site+".sql", "--quiet").Run()
}

// Correct the links with search-replace
func linkFix(surl, durl string) {
	exec.Command("wp", "search-replace", "--url="+durl, "--all-tables-with-prefix", surl, durl, "--quiet").Run()
}

// Copy the site assets over
func assetCopy(fpath, sid, tpath, did string) {
	exec.Command("rsync", "-a", "/data/www-assets/"+fpath+"/uploads/sites/"+sid+"/", "/data/www-assets/"+tpath+"/uploads/sites/"+did+"/").Run()
}

// Correct the uploads folder references
func folderRef(durl, sid, did string) {
	exec.Command("wp", "search-replace", "--url="+durl, "--all-tables-with-prefix", "app/uploads/sites/"+sid, "app/uploads/sites/"+did, "--quiet").Run()
}

// Catch any lingering http addresses
func httpFind(durl string) {
	exec.Command("wp", "search-replace", "--url="+durl, "--all-tables-with-prefix", "http://", "https://", "--quiet").Run()
}

// Remap the users to match their new ID
func remap(durl string) {
	exec.Command("/bin/bash", "-c", "/data/scripts/user_import.py", "-p", "current/web/wp", "-u", durl, "-i ", "/data/temp/"+site+".json").Run()
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
