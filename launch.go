package main

import (
	"encoding/json"
	"os"
)

// Blog holds converted JSON data
type Blog struct {
	BlogID      string `json:"blog_id"`
	URL         string `json:"url"`
	LastUpdated string `json:"last_updated"`
	Registered  string `json:"registered"`
}

// Platform holds the JSON data
type Platform struct {
	GWW        Website `json:"gww"`
	Production Website `json:"production"`
	Staging    Website `json:"staging"`
	Test       Website `json:"test"`
	Vanity     Website `json:"vanity"`
	Email      Person  `json:"email"`
}

// Website holds the JSON data
type Website struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

// Person holds the JSON data
type Person struct {
	Admin string `json:"admin"`
}

// Variable declarations
var (
	flag                                     = os.Args[1]
	websites                                 Platform
	testObj, sourceOBJ, destOBJ              Blog
	siteName, testID, stageID, prodID        string
	sourcePath, sourceURL, destPath, destURL string
)

/*
Flags:
	s2p - Staging to Production
	p2s - Production to Staging
	t2t - Test to Test
*/

// Quarterback function controls the flow of the program
func quarterback() {
	sites := readit("local/env.json")
	json.Unmarshal(sites, &websites)
	siteName = os.Args[2]

	switch flag {
	case "-s2p":
		source(websites.Staging.Path, websites.Staging.URL)
		first()
		destination(websites.Production.Path, websites.Production.URL)
		receiver()
	case "-p2s":
		source(websites.Production.Path, websites.Production.URL)
		first()
		destination(websites.Staging.Path, websites.Staging.URL)
		receiver()
	case "-t2t":
		source(websites.Vanity.Path, websites.Vanity.URL)
		first()
		destination(websites.Test.Path, websites.Test.URL)
		receiver()
	default:
		alert(huh)
	}
}

// Create the source object
func source(path, url string) {
	sourcePath, sourceURL = path, url                                       // Transfer local JSON contents to main code
	sourceList := construct(sourceURL, sourcePath)                          // List of source sites in JSON format
	sourceOBJ = aquireID("https://"+sourceURL+"/"+siteName+"/", sourceList) // Creates a specific source object
}

// Run the first few functions up to the new site creation
func first() {
	banner("Exporting the database tables")
	exportDB(sourceOBJ.URL)
	banner("Creating a user export file")
	exportUsers()
	banner("Creating the new WordPress site")
	createSite(siteName, websites.Email.Admin)
}

// Create the destination object
func destination(path, url string) {
	destPath, destURL = path, url                                     // Transfer local JSON contents to main code
	destList := construct(destURL, destPath)                          // List of destination sites in JSON format
	destOBJ = aquireID("https://"+destURL+"/"+siteName+"/", destList) // The specific destination object
}

// Query WordPress for a list of all sites and map the JSON data to a struct array
func construct(url, path string) []Blog {
	var blog []Blog
	// query, _ := exec.Command("wp", "site", "list", "--path=/data/www-app/"+path+"/current/web/wp", "--url="+url, "--format=json").Output()
	query := execute("-c", "wp", "site", "list", "--path=/data/www-app/"+path+"/current/web/wp", "--url="+url, "--format=json")
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

// Trigger the rest of the program after passing through the Quarterback
func receiver() {
	second()
	dryrun()
	last()
}

// Run the second round of functions after being able to grab the new site ID
func second() {
	banner("Backing up the database")
	backupDB()
	banner("Replacing the destination blog_id with that of the source")
	replaceIDs(sourceOBJ.BlogID, destOBJ.BlogID)
	banner("Importing the database tables")
	importDB()
}

// Pre-emptively run the data modifying functions in --dry-run mode
func dryrun() {
	banner("Updating URL's")
	direct(confirm(linkFixDR()), "lf")
	banner("Copying Assets")
	direct(confirm(assetCopyDR(sourceOBJ.BlogID, destOBJ.BlogID)), "ac")
	banner("Fixing Uploads")
	direct(confirm(uploadsFolderDR(sourceOBJ.BlogID, destOBJ.BlogID)), "fr")
	banner("Fixing Escapes")
	direct(confirm(uploadsFolderEscapesDR(sourceOBJ.BlogID, destOBJ.BlogID)), "fr2")
	banner("Fixing HTTP References")
	direct(confirm(httpFindDR()), "hf")
}

// Run the remaining functions
func last() {
	banner("Remaping the users to match their new ID")
	remap()
	banner("Flushing the WordPress cache")
	flush()
}

// Export the WordPress database tables
func exportDB(sourceURL string) {
	// exec.Command("wp", "db", "export", "--tables=$("+string(sub)+")", "--quiet", "/data/temp/"+siteName+".sql").Run()
	// exec.Command("wp", "db", "export", "--tables=$(wp db tables", "--url="+sourceURL, "--all-tables-with-prefix", "--format=csv)", "/data/temp/"+siteName+".sql").Run()
	sub := execute("-c", "wp db tables", "--url="+sourceURL+"--all-tables-with-prefix --format=csv")
	execute("-v", "wp", "db", "export", "--tables=$("+string(sub)+")", "/data/temp/"+siteName+".sql")
}

// Create a user export file in JSON format
func exportUsers() {
	// exec.Command("/bin/bash", "-c", "/data/scripts/user_export.py", "-p", "/data/www-app/"+sourcePath+"/current/web/wp", "-u", sourceURL, "-o", "/data/temp/"+siteName+".json").Run()
	// execute("-v", "/bin/bash", "-c", "/data/scripts/user_export.py", "-p", "/data/www-app/"+sourcePath+"/current/web/wp", "-u", sourceURL, "-o", "/data/temp/"+siteName+".json")
	people := execute("-c", "wp", "user", "list", "--url="+sourceURL, "--path="+"/data/www-app/"+sourcePath+"/current/web/wp", "--format=json")
	inspect(os.WriteFile("/data/temp/"+siteName+".json", people, 0666))
}

// Create a new WordPress site
func createSite(title, email string) {
	// exec.Command("wp", "site", "create", "--url=https://"+destURL+"/"+siteName+"/", "--title="+title, "--email="+email, "--quiet").Run()
	execute("-v", "wp", "site", "create", "--url=https://"+destURL+"/"+siteName+"/", "--title="+title, "--email="+email)
}

// Backup the WordPress SQL database
func backupDB() {
	// exec.Command("wp", "db", "export", "--path=/data/www-app/"+destPath+"/current/web/wp", "/data/temp/backup.sql", "--quiet").Run()
	execute("-v", "wp", "db", "export", "--path=/data/www-app/"+destPath+"/current/web/wp", "/data/temp/backup.sql")
}

// Replace the destination (did) blog_id's with that of the source (sid)
func replaceIDs(sid, did string) {
	// exec.Command("sed", "-i", "'s/wp_"+sid+"_/wp_"+did+"_/g'", "/data/temp/"+siteName+".sql").Run()
	execute("-v", "sed", "-i", "'s/wp_"+sid+"_/wp_"+did+"_/g'", "/data/temp/"+siteName+".sql")
}

// Import the WordPress SQL database tables
func importDB() {
	// exec.Command("wp", "db", "import", "/data/temp/"+siteName+".sql", "--quiet").Run()
	execute("-v", "wp", "db", "import", "/data/temp/"+siteName+".sql")
}

// Correct the links with search-replace
func linkFix() {
	// exec.Command("wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", sourceURL, destURL, "--quiet").Run()
	execute("-v", "wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", sourceURL, destURL)
}

// Copy the WordPress site assets over
func assetCopy(sid, did string) {
	// exec.Command("rsync", "-a", "/data/www-assets/"+sourcePath+"/uploads/sites/"+sid+"/", "/data/www-assets/"+destPath+"/uploads/sites/"+did+"/").Run()
	execute("-v", "rsync", "-a", "/data/www-assets/"+sourcePath+"/uploads/sites/"+sid+"/", "/data/www-assets/"+destPath+"/uploads/sites/"+did+"/")
}

// Correct the references to the uploads folder
func uploadsFolder(sid, did string) {
	// exec.Command("wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", "app/uploads/sites/"+sid, "app/uploads/sites/"+did, "--quiet").Run()
	execute("-v", "wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", "app/uploads/sites/"+sid, "app/uploads/sites/"+did)
}

// Correct any unescaped folders due to Gutenberg Blocks
func uploadsFolderEscapes(sid, did string) {
	// exec.Command("wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", "app\\/uploads\\/sites\\/"+sid, "app\\/uploads\\/sites\\/"+did, "--quiet").Run()
	execute("-v", "wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", "app\\/uploads\\/sites\\/"+sid, "app\\/uploads\\/sites\\/"+did)
}

// Catch any lingering http addresses and convert them to https
func httpFind() {
	// exec.Command("wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", "http://", "https://", "--quiet").Run()
	execute("-v", "wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", "http://", "https://")
}

// Remap the users to match their new ID
func remap() {
	// exec.Command("/bin/bash", "-c", "/data/scripts/user_import.py", "-p", "/data/www-app/"+destPath+"/current/web/wp", "-u", destURL, "-i ", "/data/temp/"+siteName+".json").Run()
	execute("-v", "/bin/bash", "-c", "/data/scripts/user_import.py", "-p", "/data/www-app/"+destPath+"/current/web/wp", "-u", destURL, "-i ", "/data/temp/"+siteName+".json")
}

// Flush the WordPress cache
func flush() {
	// exec.Command("wp", "cache", "flush", "--quiet").Run()
	execute("-v", "wp", "cache", "flush")
}
