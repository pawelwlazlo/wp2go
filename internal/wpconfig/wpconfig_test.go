package wpconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReplaceDefine(t *testing.T) {
	// Sample matching real wp-config.php format
	input := `define( 'DB_NAME', 'i9930109_ety31' );
define( 'DB_USER', 'i9930109_ety31' );
define( 'DB_PASSWORD', 'T.nu53JtdHdvYg3Jsqm06' );
define( 'DB_HOST', 'localhost' );
$table_prefix = 'pzfx_';`

	got := replaceDefine(input, "DB_NAME", DDEVDBName)
	got = replaceDefine(got, "DB_USER", DDEVDBUser)
	got = replaceDefine(got, "DB_PASSWORD", DDEVDBPass)
	got = replaceDefine(got, "DB_HOST", DDEVDBHost)

	for _, want := range []string{
		"define( 'DB_NAME', 'db' )",
		"define( 'DB_USER', 'db' )",
		"define( 'DB_PASSWORD', 'db' )",
		"define( 'DB_HOST', 'db' )",
		"$table_prefix = 'pzfx_';",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output should contain %q;\ngot:\n%s", want, got)
		}
	}
	if strings.Contains(got, "i9930109_ety31") || strings.Contains(got, "localhost") || strings.Contains(got, "T.nu53") {
		t.Errorf("old DB values should be replaced; got:\n%s", got)
	}
}

func TestPatchForDDEV(t *testing.T) {
	dir := t.TempDir()
	wpConfig := `<?php
define( 'DB_NAME', 'old_db' );
define( 'DB_USER', 'old_user' );
define( 'DB_PASSWORD', 'old_pass' );
define( 'DB_HOST', 'localhost' );
$table_prefix = 'wp_';
`
	if err := os.WriteFile(filepath.Join(dir, "wp-config.php"), []byte(wpConfig), 0644); err != nil {
		t.Fatal(err)
	}

	if err := PatchForDDEV(dir); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "wp-config.php"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for _, s := range []string{"'db'", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_HOST", "$table_prefix = 'wp_';"} {
		if !strings.Contains(content, s) {
			t.Errorf("patched file should contain %q", s)
		}
	}
	if strings.Contains(content, "old_db") || strings.Contains(content, "old_user") || strings.Contains(content, "localhost") {
		t.Errorf("old credentials should be replaced: %s", content)
	}
}
