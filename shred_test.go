package main 

import (
	"testing"
	"os"
	"fmt"
)

func createFile(t *testing.T, name string, data []byte) string {
	err := os.WriteFile(name, data, 0644)
	if err != nil {
		t.Fatal("Failed to create temp file: ", err)
	}
	return name
}

// Test the core functionality of the program 
func TestShred(t *testing.T){
	fmt.Println("TestShred")
	// Create file 
	fileName := createFile(t, "randomfile.txt", []byte("Hello, Canonical!"))
	
	// Call shred function 
	err := Shred(fileName, 3)
	if err != nil {
		t.Fatal("Shred failed: ", err)
	}

	// Verify file is deleted
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		t.Error("File was not deleted")
	}
}

// Test shredding an empty file 
func TestShredEmptyFile(t *testing.T) {
	fmt.Println("\nTestShredEmptyFile")
	// Create empty file 
	fileName := createFile(t, "emptyfile.txt", []byte{})


	// Call shred function 
	err := Shred(fileName, 3)
	if err != nil {
		t.Fatal("Shred failed on empty file: ", err)
	}

	// Test if deleted
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		t.Error("Empty file was not deleted")
	}
}

// Test if file does not exist 
func TestShredFileDoesNotExist(t *testing.T) {
	fmt.Println("\nTestShredFileDoesNotExist")
	err := Shred("unknown.txt", 3)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

// Stress test 
func TestShredLargeFile(t *testing.T) {
	fmt.Println("\nTestShredLargeFile")
	// Create ~10MB file for testing 
	data := make([]byte, 10*1024*1024)
	file := createFile(t, "big.bin", data)

	if err := Shred(file, 3); err != nil {
		t.Fatal("Shred failed on big file: ", err)
	}

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Error("Big file was not deleted")
	}
}


// Test if file is in use
func TestShredFileLocked(t *testing.T) {
	fmt.Println("\nTestShredFileLocked")
	file := createFile(t, "locked.txt", []byte("file locked"))

	// Open file
	f, err := os.Open(file)
	if err != nil {
		t.Fatal("Failed to open file: ", err)
	}
	defer f.Close()

	// Shred and log error if it occurs
	if err := Shred(file, 3); err != nil {
		t.Logf("Shred returned an error (may happen if file is locked): %v", err)
	}

	// Clean up manually
	f.Close()
	os.Remove(file)
}

// Test if wrong permission
func TestShredPermissionDenied(t *testing.T) {
	fmt.Println("\nTestShredPermissionDenied")
	file := createFile(t, "readonly.txt", []byte("I am a read-only file."))

	// Make file read-only
	if err := os.Chmod(file, 0444); err != nil {
		t.Fatal("Failed to change file permissions: ", err)
	}

	// Attempt shred 
	if err := Shred(file, 3); err != nil {
		t.Fatal("Shred failed due to permissions test: ", err)
	}
	// Reset permissions 
	os.Chmod(file, 0644)
	os.Remove(file)

}

// Shred filenames with special characters 
func TestShredSpecialFilenames(t *testing.T) {
	fmt.Println("\nTestShredSpecialFilenames")
	file := createFile(t, "spécial file @#$.txt", []byte("data"))
	defer os.Remove(file)

	if err := Shred(file, 3); err != nil {
		t.Errorf("Shred failed on special filename: %v", err)
	}
}

// Test shred hidden files 
func TestShredHiddenFile(t *testing.T) {
	fmt.Println("\nTestShredHiddenFile")
	file := createFile(t, ".hiddenfile", []byte("hidden"))
	defer os.Remove(file)

	if err := Shred(file, 3); err != nil {
		t.Errorf("Shred failed on hidden file: %v", err)
	}
}

// Shred file twice at the same time
func TestShredConcurrent(t *testing.T) {
	fmt.Println("\nTestShredConcurrent")
	file := createFile(t, "concurrent.txt", []byte("data"))
	defer os.Remove(file)

	done := make(chan bool)
	for i := 0; i < 2; i++ {
		go func() {
			_ = Shred(file, 3) // ignore error; just ensure no panic
			done <- true
		}()
	}
	<-done
	<-done
}

// Test function for large filename 
func TestShredLongFilename(t *testing.T) {
	fmt.Println("\nTestShredLongFilename")
	name := "l"
	for i := 0; i < 250; i++ {
		name += "l"
	}
	name += ".txt"
	file := createFile(t, name, []byte("long filename"))
	defer os.Remove(file)

	if err := Shred(file, 3); err != nil {
		t.Errorf("Shred failed on long filename: %v", err)
	}
}