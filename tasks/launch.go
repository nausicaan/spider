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
	flag     = os.Args[1]
	siteName = os.Args[2]

	testID, stageID, prodID                  string
	testObj, sourceOBJ, destOBJ              Blog
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
		sourcePath, sourceURL = stagePath, stageURL                            //transfer secret constants to main code
		stageList = parseJSON(sourceURL, sourcePath)                           // List of Stagging sites in JSON format
		sourceOBJ = aquireID("https://"+sourceURL+"/"+siteName+"/", stageList) // Creates a specific Stagging object
		first()
		destPath, destURL = prodPath, prodURL                             //transfer secret constants to main code
		prodList = parseJSON(destURL, destPath)                           // List of Production sites in JSON format
		destOBJ = aquireID("https://"+destURL+"/"+siteName+"/", prodList) // The specific Production object
	case "p2s":
		sourcePath, sourceURL = prodPath, prodURL                             //transfer secret constants to main code
		prodList = parseJSON(sourceURL, sourcePath)                           // List of Production sites in JSON format
		sourceOBJ = aquireID("https://"+sourceURL+"/"+siteName+"/", prodList) // The specific Production object
		first()
		destPath, destURL = stagePath, stageURL                            //transfer secret constants to main code
		stageList = parseJSON(destURL, destPath)                           // List of Stagging sites in JSON format
		destOBJ = aquireID("https://"+destURL+"/"+siteName+"/", stageList) // Creates a specific Stagging object
	default:
		testList = parseJSON(testURL, testPath)
		sourceOBJ = aquireID("http://test.engage.gov.bc.ca/"+siteName+"/", testList)
		first()
	}
	second()
	dryruns()
	last()
}

// Run the first few functions up to the new site creation
func first() {
	exportDB(sourceOBJ.URL)
	exportUsers()
	createSite(siteName, adminEmail)
}

// Run the remaining functions after being able to grab the new site ID
func second() {
	backupDB()
	replaceIDs(sourceOBJ.BlogID, destOBJ.BlogID)
	importDB()
}

func dryruns() {
	title("|U| |P| |D| |A| |T| |E| | | |U| |R| |L| |S|")
	direct(confirm(linkFixDR()), "lf")
	title("|C| |O| |P| |Y| | | |A| |S| |S| |E| |T| |S|")
	direct(confirm(assetCopyDR(sourceOBJ.BlogID, destOBJ.BlogID)), "ac")
	title("|F| |I| |X| | | |U| |P| |L| |O| |A| |D| |S|")
	direct(confirm(uploadsFolderDR(sourceOBJ.BlogID, destOBJ.BlogID)), "fr")
	title("|F| |I| |X| | | |E| |S| |C| |A| |P| |E| |S|")
	direct(confirm(uploadsFolderEscapesDR(sourceOBJ.BlogID, destOBJ.BlogID)), "fr2")
	title("|F| |I| |X| | | |H| |T| |T| |P| |:| |/| |/|")
	direct(confirm(httpFindDR()), "hf")
}

func last() {
	remap()
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
func exportDB(sourceURL string) {
	sub, err := exec.Command("wp db tables", "--url="+sourceURL+"--all-tables-with-prefix --format=csv").Output()
	errors(err)
	exec.Command("wp", "db", "export", "--tables=$("+string(sub)+")", "--quiet", "/data/temp/"+siteName+".sql").Run()
	// exec.Command("wp", "db", "export", "--tables=$(wp db tables", "--url="+sourceURL, "--all-tables-with-prefix", "--format=csv)", "/data/temp/"+siteName+".sql").Run()
}

// Create a user export file
func exportUsers() {
	exec.Command("/bin/bash", "-c", "/data/scripts/user_export.py", "-p", "/data/www-app/"+sourcePath+"/current/web/wp", "-u", sourceURL, "-o", "/data/temp/"+siteName+".json").Run()
}

// Create the new WordPress site
func createSite(title, email string) {
	exec.Command("wp", "site", "create", "--url=https://"+destURL+"/"+siteName+"/", "--title="+title, "--email="+email, "--quiet").Run()
}

// Backup the database
func backupDB() {
	exec.Command("wp", "db", "export", "--path=/data/www-app/"+destPath+"/current/web/wp", "/data/temp/backup.sql", "--quiet").Run()
}

// Take the blog_id from the source (sid) and send it to the destination (did) to be replaced
func replaceIDs(sid, did string) {
	exec.Command("sed", "-i", "'s/wp_"+sid+"_/wp_"+did+"_/g'", "/data/temp/"+siteName+".sql").Run()
}

// Import the data
func importDB() {
	exec.Command("wp", "db", "import", "/data/temp/"+siteName+".sql", "--quiet").Run()
}

// Correct the links with search-replace
func linkFix() {
	exec.Command("wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", sourceURL, destURL, "--quiet").Run()
}

// Copy the site assets over
func assetCopy(sid, did string) {
	exec.Command("rsync", "-a", "/data/www-assets/"+sourcePath+"/uploads/sites/"+sid+"/", "/data/www-assets/"+destPath+"/uploads/sites/"+did+"/").Run()
}

// Correct the references to the uploads folder
func uploadsFolder(sid, did string) {
	exec.Command("wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", "app/uploads/sites/"+sid, "app/uploads/sites/"+did, "--quiet").Run()
}

// Correct any unescaped folders due to Gutenberg Blocks
func uploadsFolderEscapes(sid, did string) {
	exec.Command("wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", "app\\/uploads\\/sites\\/"+sid, "app\\/uploads\\/sites\\/"+did, "--quiet").Run()
}

// Catch any lingering http addresses
func httpFind() {
	exec.Command("wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", "http://", "https://", "--quiet").Run()
}

// Remap the users to match their new ID
func remap() {
	exec.Command("/bin/bash", "-c", "/data/scripts/user_import.py", "-p", "/data/www-app/"+destPath+"/current/web/wp", "-u", destURL, "-i ", "/data/temp/"+siteName+".json").Run()
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
