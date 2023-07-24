
# randnumsplash
A small and fast package that generates a dummy file of the desired length containing random numbers."

## Installation

```bash
go get -u github.com/anurag0608/randnumsplash
```

## Usage

```golang
import "github.com/anurag0608/randnumsplash"
// GenerateRandFile(targetSizeInBytes, folder_location, fileName, loggingEnabled)
err := randnumsplash.GenerateRandFile(1024*1024, folder_location, "rand.txt", true) // loggingEnable
```
Sample output (loggingEnabled = true)
```bash
Target file size: 535822336 Bytes | 511.00 Mb
Starting dumping random numbers 🤖
📂 100% |███████████████| (511/511 MB, 21 MB/s)
Done✅
Took ✨ 24.31s
```
## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

Apache License