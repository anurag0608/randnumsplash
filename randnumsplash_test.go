// write test boiler code
package randnumsplash

import (
	"os"
	"testing"
)

func TestGenerateRandFile(t *testing.T) {
	defer func(){
		os.Remove("testGenerateRandFile")
	}()
	// generate test file
	var targetFileSize int64 = 1024*1024 // 1Mb
	err := GenerateRandFile(targetFileSize, "", "testGenerateRandFile", false)
	if err!=nil{
		t.Fatal(err)
	}
}