package sql

import (
	"os"
	"regexp"
	"strings"
)

// Prepare reads the SQL file, detects the old site URL from wp_options (siteurl/home),
// replaces it with newFQDN, writes to a temp file, and returns the temp file path.
// Caller must call Cleanup with the returned path when done.
func Prepare(sqlPath, newFQDN string) (tempPath string, err error) {
	data, err := os.ReadFile(sqlPath)
	if err != nil {
		return "", err
	}

	oldURL, err := detectOldURL(data)
	if err != nil || oldURL == "" {
		oldURL = "http://localhost"
	}

	content := string(data)
	content = replaceURL(content, oldURL, newFQDN)

	f, err := os.CreateTemp("", "wp2go-*.sql")
	if err != nil {
		return "", err
	}
	tempPath = f.Name()
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		os.Remove(tempPath)
		return "", err
	}
	if err := f.Close(); err != nil {
		os.Remove(tempPath)
		return "", err
	}
	return tempPath, nil
}

// Cleanup removes the temp file created by Prepare.
func Cleanup(tempPath string) {
	os.Remove(tempPath)
}

var siteurlRe = regexp.MustCompile(`'siteurl','(https?://[^']+)'`)
var homeRe = regexp.MustCompile(`'home','(https?://[^']+)'`)

func detectOldURL(data []byte) (string, error) {
	s := string(data)
	for _, re := range []*regexp.Regexp{siteurlRe, homeRe} {
		if m := re.FindStringSubmatch(s); len(m) > 1 {
			return m[1], nil
		}
	}
	return "", nil
}

func replaceURL(content, oldURL, newFQDN string) string {
	content = strings.ReplaceAll(content, oldURL, newFQDN)
	// Replace the other scheme variant (e.g. if oldURL was http://, also replace https://)
	host := strings.TrimPrefix(oldURL, "http://")
	host = strings.TrimPrefix(host, "https://")
	if host != oldURL {
		other := "https://" + host
		if strings.HasPrefix(oldURL, "https://") {
			other = "http://" + host
		}
		content = strings.ReplaceAll(content, other, newFQDN)
	}
	return content
}
