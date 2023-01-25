package tasks

import "os/exec"

// ### Dry Run Functions ###

func linkFixDR() string {
	dr, _ := exec.Command("wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", sourceURL, destURL, "--dry-run").Output()
	return string(dr)
}

func assetCopyDR(sourcePath, fid, destPath, tid string) string {
	dr, _ := exec.Command("rsync", "-a", "/data/www-assets/"+sourcePath+"/uploads/sites/"+fid+"/", "/data/www-assets/"+destPath+"/uploads/sites/"+tid+"/", "--stats", "--dry-run").Output()
	return string(dr)
}

func httpFindDR() []byte {
	dr, _ := exec.Command("wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", "http://", "https://", "--dry-run").Output()
	return dr
}

// Correct the uploads folder references
func folderRefDR(fid, tid string) string {
	dr, _ := exec.Command("wp", "search-replace", "--url="+destURL+"/"+siteName, "--all-tables-with-prefix", "app/uploads/sites/"+fid, "app/uploads/sites/"+tid, "--dry-run").Output()
	return string(dr)
}
