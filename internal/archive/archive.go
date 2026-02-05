package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extract decompresses the .tar.gz archive into projectDir. If the archive
// has a single top-level directory, its contents are written to projectDir root.
func Extract(archivePath, projectDir string) error {
	singleRoot, err := detectSingleRoot(archivePath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return err
	}

	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		destName := h.Name
		if singleRoot != "" {
			if h.Name == singleRoot {
				continue
			}
			if strings.HasPrefix(h.Name, singleRoot+"/") {
				destName = strings.TrimPrefix(h.Name, singleRoot+"/")
			}
		}
		dest := filepath.Join(projectDir, filepath.FromSlash(destName))
		switch h.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(dest, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
				return err
			}
			out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(h.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		}
	}
	return nil
}

func detectSingleRoot(archivePath string) (string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	var singleRoot string
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		name := filepath.Clean(h.Name)
		if name == "." || name == ".." || strings.HasPrefix(name, "..") {
			continue
		}
		parts := strings.SplitN(name, "/", 2)
		top := parts[0]
		if singleRoot == "" {
			singleRoot = top
		} else if top != singleRoot {
			singleRoot = ""
			break
		}
	}
	return singleRoot, nil
}
