// write test boiler code
package randnumsplash

import (
	"fmt"
	"math/rand"
	"os"
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
