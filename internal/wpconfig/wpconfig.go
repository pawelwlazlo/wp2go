package wpconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// DDEV database credentials (standard for DDEV WordPress).
const (
	DDEVDBName = "db"
	DDEVDBUser = "db"
	DDEVDBPass = "db"
	DDEVDBHost = "db"
)

// PatchForDDEV updates wp-config.php in projectDir so it uses DDEV's database
// credentials. Existing define('DB_NAME', ...) (and DB_USER, DB_PASSWORD,
// DB_HOST) are replaced; $table_prefix and the rest of the file are left unchanged.
func PatchForDDEV(projectDir string) error {
	path := filepath.Join(projectDir, "wp-config.php")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read wp-config.php: %w", err)
	}

	content := string(data)
	content = replaceDefine(content, "DB_NAME", DDEVDBName)
	content = replaceDefine(content, "DB_USER", DDEVDBUser)
	content = replaceDefine(content, "DB_PASSWORD", DDEVDBPass)
	content = replaceDefine(content, "DB_HOST", DDEVDBHost)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("write wp-config.php: %w", err)
	}
	return nil
}

// replaceDefine replaces define('NAME', 'oldvalue') or define("NAME", "oldvalue")
// with define( 'NAME', 'newvalue' ). Handles optional whitespace.
func replaceDefine(content, name, newValue string) string {
	// Match: define ( optional spaces ) ( optional spaces ) ' or " NAME ' or " , optional spaces ' or " anything ' or " optional spaces ) ;
	pattern := regexp.MustCompile(
		`define\s*\(\s*['"]` + regexp.QuoteMeta(name) + `['"]\s*,\s*['"][^'"]*['"]\s*\)`,
	)
	replacement := "define( '" + name + "', '" + newValue + "' )"
	return pattern.ReplaceAllString(content, replacement)
}
