package syncfiles

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSync(t *testing.T) {
	parentFolder, err := os.MkdirTemp("", "sync_test")
	if err != nil {
		t.Fatalf("Failed to create temporary folder: %v", err)
	}
	defer os.RemoveAll(parentFolder)

	update := `
//// FILE~test/test1.txt ////
Hello, World!
//// END FILE ////
`

	if err := Update(parentFolder, update); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(parentFolder, "test/test1.txt"))
	if err != nil {
		t.Fatalf("Failed to read synced file: %v", err)
	}

	expectedContent := "Hello, World!"
	if string(content) != expectedContent {
		t.Fatalf("Unexpected file content. Got: %q, Expected: %q", content, expectedContent)
	}
}

func TestLoad(t *testing.T) {
	parentFolder, err := os.MkdirTemp("", "load_test")
	if err != nil {
		t.Fatalf("Failed to create temporary folder: %v", err)
	}
	defer os.RemoveAll(parentFolder)

	testFolderPath := filepath.Join(parentFolder, "test")
	if err := os.Mkdir(testFolderPath, os.ModePerm); err != nil {
		t.Fatalf("Failed to create test folder: %v", err)
	}

	testFilePath := filepath.Join(testFolderPath, "test1.txt")
	testContent := []byte("Hello, World!")
	if err := os.WriteFile(testFilePath, testContent, os.ModePerm); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	loaded, err := Load(parentFolder)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	expectedLoaded := "//// FILE~test/test1.txt ////\nHello, World!\n//// END FILE ////"
	if loaded != expectedLoaded {
		t.Fatalf("Unexpected loaded content. Got: %q, Expected: %q", loaded, expectedLoaded)
	}
}
