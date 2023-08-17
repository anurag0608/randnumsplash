// write test boiler code
package randnumsplash

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestGenerateRandFile(t *testing.T) {
	defer func() {
		os.Remove("testGenerateRandFile")
	}()
	// generate test file
	var targetFileSize int64 = 1024 * 1024 // 1Mb
	err := GenerateRandFile(targetFileSize, "", "testGenerateRandFile", false)
	err2 := GenerateRandFile(targetFileSize, "", "testGenerateRandFile", true)

	if err != nil || err2 != nil {
		t.Fatal(fmt.Errorf("Error1: %w, \nError2: %w", err, err2))
	}
}
func TestInvalidTargetFileLocation(t *testing.T) {
	const length = 10
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	err := GenerateRandFile(1024*1024, string(result), "myfile", false)
	if err == nil {
		t.Fatal(err)
	}
}

func TestGenerateRandFileSynchronization(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testfiles")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	fmt.Println(tempDir)
	defer os.RemoveAll(tempDir)
	i := 0
	counter := 0
	var wg sync.WaitGroup
	// each iteration will spawn two threads
	// Total expected created files: 5*2=10
	for i < 5 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			// test file: rand_test_i%d_c%d.txt, example: rand_test_i4_c0.txt
			GenerateRandFile(1024*1024*8, tempDir, fmt.Sprintf("rand_test_i%d_c%d.txt", i, counter), false)
		}()
		go func() {
			defer wg.Done()
			GenerateRandFile(1024*1024*8, tempDir, fmt.Sprintf("rand_test_i%d_c%d.txt", i, counter+1), false)
		}()
		// wait for go routines to finish
		wg.Wait()
		i++
		counter++
	}
	// verify that all files are generated and has size approx 8Mb(+/-512kb, due to variable buffer length)
	i = 0
	checkCounter := 0
	for i < 5 {
		filePath := filepath.Join(tempDir, fmt.Sprintf("rand_test_i%d_c%d.txt", i, checkCounter))
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("failed to get file stats: %v", err)
		}
		if ok, _ := isApproximatelyXMB(filePath, float64(8)); !ok {
			t.Errorf("file size mismatch: expected ~1024*1024*8(+/-512Kb allowed), got %d", fileInfo.Size())
		}
		i++
		checkCounter++
	}
}
func isApproximatelyXMB(filePath string, X float64) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to get file stats: %w", err)
	}
	fileSize := fileInfo.Size()
	return float64(fileSize) >= 1024*1024*X-512*1024 && float64(fileSize) <= 1024*1024*X+512*1024, nil
}
