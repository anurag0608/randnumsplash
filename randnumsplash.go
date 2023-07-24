package randnumsplash

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/schollz/progressbar/v3"
)

var random *rand.Rand

func init() {
	// Seed the random number generator
	seed := time.Now().UnixNano() + int64(time.Now().Nanosecond()) + int64(rand.Intn(1000))
	random = rand.New(rand.NewSource(seed))
}
func genRandNum() int64 {
	// Generate a random int64 number between 0 and the maximum 10000
	return random.Int63n(10000)
}
func GenerateRandFile(targetFileSizeInBytes int64, targetLocation, fileName string, loggingEnabled bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("PANIC: %v", r)
		}
	}()
	var file *os.File
	targetFileLoc := filepath.Join(targetLocation, fileName)
	if _, err := os.Stat(targetFileLoc); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to get file stats: %w", err)
		}
		// file don't exist
		if loggingEnabled {
			fmt.Println("File:" + fileName + " not found ⚠️\nCreating... ✅")
		}
		file, err = os.Create(targetFileLoc)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	} else {
		// file exist
		file, err = os.OpenFile(targetFileLoc, os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
	}
	defer file.Close()

	// get sleep timer for updateFileSizeRoutine and buffer size using targetFileSize
	buffSize := computeBufferSize(targetFileSizeInBytes)

	// create a buffered writer, with buffer of size 64kb
	bw := bufio.NewWriterSize(file, buffSize)
	if loggingEnabled {
		fmt.Printf("Target file size: %v Bytes | %0.2f Mb\n", targetFileSizeInBytes, float64(targetFileSizeInBytes)/(1024*1024))
		fmt.Println("Starting dumping random numbers 🤖")
	}
	// progressLineItr := utils.GetProgressLineIterator()
	start := time.Now()

	progressBar := progressbar.NewOptions64(
		targetFileSizeInBytes,
		progressbar.OptionSetDescription(`📂`),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetRenderBlankState(false),
	)
	var total int64 = 0
	for {
		if total >= targetFileSizeInBytes {
			break
		}
		str := fmt.Sprintf("%v\n", genRandNum())
		total += int64(len([]byte(str)))
		if loggingEnabled {
			progressBar.Add(len([]byte(str)))
		}
		bw.WriteString(str)
	}
	if loggingEnabled {
		progressBar.Finish()
	}
	// dont flush current file size already exceeded the target size
	// bw.Flush()
	end := time.Now()
	elapsed := end.Sub(start)
	if loggingEnabled {
		fmt.Printf("Done✅ \nTook ✨ %0.2fs\n", elapsed.Seconds())
	}
	return nil
}
func computeBufferSize(targetFileSizeInBytes int64) int {
	buffSize := 64 * 1024
	if targetFileSizeInBytes < 64*1024 {
		buffSize = 1
	}
	return buffSize
}
