package domain

import "strings"

// Normalize returns the short name (for directory and DDEV project) and the
// full URL (FQDN) for the site. If sitename has no dot, FQDN is https://<sitename>.ddev.site.
func Normalize(sitename string) (shortName, fqdn string) {
	sitename = strings.TrimSpace(sitename)
	if sitename == "" {
		return "", ""
	}
	if strings.Contains(sitename, ".") {
		shortName = strings.TrimSuffix(sitename, ".ddev.site")
		if shortName == sitename {
			shortName = strings.Split(sitename, ".")[0]
		}
		if strings.HasPrefix(sitename, "http://") || strings.HasPrefix(sitename, "https://") {
			fqdn = sitename
		} else {
			fqdn = "https://" + sitename
		}
		return shortName, fqdn
	}
	return sitename, "https://" + sitename + ".ddev.site"
}
