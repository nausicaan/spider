package tasks

import "os/exec"

// ### Dry Run Functions ###

func linkFixDR(furl, turl string) string {
	dr, _ := exec.Command("wp", "search-replace", "--url="+turl+"/"+site, "--all-tables-with-prefix", furl, turl, "--dry-run").Output()
	return string(dr)
}

func assetCopyDR(fpath, fid, tpath, tid string) string {
	dr, _ := exec.Command("rsync", "-a", "/data/www-assets/"+fpath+"/uploads/sites/"+fid+"/", "/data/www-assets/"+tpath+"/uploads/sites/"+tid+"/", "--stats", "--dry-run").Output()
	return string(dr)
}

func httpFindDR(turl string) []byte {
	dr, _ := exec.Command("wp", "search-replace", "--url="+turl+"/"+site, "--all-tables-with-prefix", "http://", "https://", "--dry-run").Output()
	return dr
}

// Correct the uploads folder references
func folderRefDR(turl, fid, tid string) string {
	dr, _ := exec.Command("wp", "search-replace", "--url="+turl+"/"+site, "--all-tables-with-prefix", "app/uploads/sites/"+fid, "app/uploads/sites/"+tid, "--dry-run").Output()
	return string(dr)
}
