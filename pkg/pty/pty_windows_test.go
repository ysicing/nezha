//go:build windows && !arm64

package pty

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestSecureUnzip_PathTraversal(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		wantErr  bool
	}{
		{"Normal file", "test.txt", false},
		{"Subdirectory file", "subdir/test.txt", false},
		{"Path traversal with ..", "../evil.txt", true},
		{"Path traversal nested", "../../evil.txt", true},
		{"Absolute path Unix", "/etc/passwd", true},
		{"Absolute path Windows", "C:/Windows/System32/evil.exe", true},
		{"Hidden traversal", "subdir/../../../evil.txt", true},
		{"Current directory", ".", true},
		{"Empty path", "", true},
	}

	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "unzip_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test zip file
			zipPath := filepath.Join(tempDir, "test.zip")
			if err := createTestZip(zipPath, tt.fileName); err != nil {
				t.Fatal(err)
			}
			defer os.Remove(zipPath)

			// Test unzip
			destDir := filepath.Join(tempDir, "dest")
			err := secureUnzip(zipPath, destDir)

			if tt.wantErr && err == nil {
				t.Errorf("secureUnzip() expected error for %s, but got none", tt.fileName)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("secureUnzip() unexpected error for %s: %v", tt.fileName, err)
			}
		})
	}
}

func createTestZip(zipPath, fileName string) error {
	file, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	f, err := w.Create(fileName)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte("test content"))
	return err
}
