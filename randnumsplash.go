package randnumsplash

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/anurag0608/randnumsplash/utils"
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
		fmt.Println("File:" + fileName + " not found ⚠️\nCreating... ✅")
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
	sleep, buffSize := computeSleepAndBuffer(targetFileSizeInBytes)

	// create a buffered writer, with buffer of size 64kb
	bw := bufio.NewWriterSize(file, buffSize)
	if loggingEnabled {
		fmt.Printf("Target file size: %v Bytes | %0.2f Mb\n", targetFileSizeInBytes, float64(targetFileSizeInBytes)/(1024*1024))
		fmt.Println("starting dumping random numbers 🤖")
	}
	progressLineItr := utils.GetProgressLineIterator()
	start := time.Now()
	var currentFileSize int64 = 0
	done := make(chan bool)
	// start go routine
	go updateFileSizeRoutine(targetFileLoc, done, &currentFileSize, sleep)

	for {
		if currentFileSize >= targetFileSizeInBytes {
			break
		}
		if loggingEnabled {
			utils.ShowProgressBar(currentFileSize, targetFileSizeInBytes, progressLineItr)
		}
		bw.WriteString(fmt.Sprintf("%v\n", genRandNum()))
	}
	// exit the updateFileSizeRoutine
	close(done)
	// show remaining progress bar
	if loggingEnabled {
		utils.ShowProgressBar(currentFileSize, targetFileSizeInBytes, progressLineItr)
	}
	// dont flush current file size already exceeded the target size
	// bw.Flush()
	end := time.Now()
	elapsed := end.Sub(start)
	if loggingEnabled {
		fmt.Printf("\nDone✅ \nTook ✨ %0.2fs\n", elapsed.Seconds())
	}
	return nil
}

// go routine for updating current file size, to avoid overhead
func updateFileSizeRoutine(fileName string, done chan bool, currentFileSize *int64, sleep time.Duration) {
	for {
		select {
		case <-done:
			// exit the goroutine
			return
		default:
			fileinfo, err := os.Stat(fileName)
			if err == nil {
				*currentFileSize = fileinfo.Size()
			}
			time.Sleep(sleep)
		}
	}
}
func computeSleepAndBuffer(targetFileSizeInBytes int64) (time.Duration, int) {
	var sleep time.Duration
	buffSize := 64 * 1024
	if targetFileSizeInBytes < 64*1024 {
		sleep = 0
		buffSize = 1
	} else {
		sleep = time.Second
	}
	return sleep, buffSize
}
