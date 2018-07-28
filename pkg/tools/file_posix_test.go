// +build linux darwin

package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileMode(t *testing.T) {
	type fileModeTest struct {
		filename string
		perm     string
	}

	fileModeTests := []fileModeTest{
		{"perm-0400.txt", "0400"},
		{"perm-0644.txt", "0644"},
		{"perm-0665.txt", "0665"},
		{"perm-0755.txt", "0755"},
		{"perm-0777.txt", "0777"},
		{"perm-1644.txt", "1644"},
		{"perm-2644.txt", "2644"},
		{"perm-4655.txt", "4655"},
	}

	testDataDir, err := ensureTestDataDir()
	if err != nil {
		t.Fatalf("failed to create test_data: %v", err)
	}

	err = createTestFiles(testDataDir)
	if err != nil {
		t.Fatalf("failed to create test files: %v", err)
	}

	for _, fm := range fileModeTests {
		path := filepath.Join(testDataDir, fm.filename)
		fi, err := os.Lstat(path)
		if err != nil {
			t.Errorf("failed to get file info: %v", err)
		}

		actual := FileMode(fi.Mode())
		if actual != fm.perm {
			t.Errorf("expected %s, got %s", fm.perm, actual)
		}
	}

	err = destroyTestFiles(testDataDir)
	if err != nil {
		t.Fatalf("failed to cleanup test_data: %v", err)
	}
}
