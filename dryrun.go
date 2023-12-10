package main

/* ----- Dry Run Functions ----- */

// Correct the links with search-replace using --dry-run
func linkFixDR() string {
	// dr, _ := exec.Command("wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", sourceURL, destURL, "--dry-run").Output()
	dr := execute("-c", "wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", sourceURL, destURL, "--dry-run")
	return string(dr)
}

// Copy the site assets over using --dry-run
func assetCopyDR(sid, did string) string {
	// dr, _ := exec.Command("rsync", "-a", "/data/www-assets/"+sourcePath+"/uploads/sites/"+sid+"/", "/data/www-assets/"+destPath+"/uploads/sites/"+did+"/", "--stats", "--dry-run").Output()
	dr := execute("-c", "rsync", "-a", "/data/www-assets/"+sourcePath+"/uploads/sites/"+sid+"/", "/data/www-assets/"+destPath+"/uploads/sites/"+did+"/", "--stats", "--dry-run")
	return string(dr)
}

// Correct the references to the uploads folder using --dry-run
func uploadsFolderDR(sid, did string) string {
	// dr, _ := exec.Command("wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", "app/uploads/sites/"+sid, "app/uploads/sites/"+did, "--dry-run").Output()
	dr := execute("-c", "wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", "app/uploads/sites/"+sid, "app/uploads/sites/"+did, "--dry-run")
	return string(dr)
}

// Correct any unescaped folders due to Gutenberg Blocks using --dry-run
func uploadsFolderEscapesDR(sid, did string) string {
	// dr, _ := exec.Command("wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", "app\\/uploads\\/sites\\/"+sid, "app\\/uploads\\/sites\\/"+did, "--dry-run").Output()
	dr := execute("-c", "wp", "search-replace", "--url="+destURL, "--all-tables-with-prefix", "app\\/uploads\\/sites\\/"+sid, "app\\/uploads\\/sites\\/"+did, "--dry-run")
	return string(dr)
}

// Catch any lingering http addresses using --dry-run
func httpFindDR() string {
	// dr, _ := exec.Command("wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", "http://", "https://", "--dry-run").Output()
	dr := execute("-c", "wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", "http://", "https://", "--dry-run")
	return string(dr)
}
