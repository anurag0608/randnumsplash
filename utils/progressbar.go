package utils

import "fmt"

/*Too Slow: Not using it anymore*/
func ShowProgressBar(currentFileSizeInBytes int64, targetFileSizeInBytes int64, progressLineItr func() string) {
	fmt.Printf(
		"\r%v %v/%v Bytes | %v/%v Kb | %0.2f/%0.2f Mb",
		progressLineItr(),
		currentFileSizeInBytes, targetFileSizeInBytes,
		currentFileSizeInBytes/1024, targetFileSizeInBytes/1024,
		float64(currentFileSizeInBytes)/(1024*1024), float64(targetFileSizeInBytes)/(1024*1024),
	)
}
func GetProgressLineIterator() func() string {
	lines := []string{"|", "|", "|", "|", "/", "/", "/", "/", "-", "-", "-", "-", "\\", "\\", "\\", "\\"}
	i := -1
	return func() string {
		i = (i + 1) % len(lines)
		return lines[i]
	}
}
