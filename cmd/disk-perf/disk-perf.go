package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	bytes    = flag.String("bytes", "1g", "the number of bytes to read")
	filename = flag.String("file", "", "the file to write to")
)

const (
	KB = 1024 << (10 * iota)
	MB
	GB
)

func main() {
	flag.Parse()

	var err error
	var file *os.File
	if *filename == "" {
		file, err = os.CreateTemp("", "")
		if err == nil {
			defer os.Remove(file.Name())
		}
	} else {
		file, err = os.Create(*filename)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	start := time.Now()
	written, err := run(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	duration := time.Since(start)
	fmt.Printf("Wrote %d bytes in %s (%s)\n", written, duration, humanByteRate(float64(written)/duration.Seconds()))
}

func parseHumanBytes(bytes string) (res int, err error) {
	bytes = strings.ToLower(bytes)
	factor := 1
	switch {
	case strings.HasSuffix(bytes, "k"):
		factor = 1024
		bytes = bytes[0 : len(bytes)-1]
	case strings.HasSuffix(bytes, "m"):
		factor = 1024 << 10
		bytes = bytes[0 : len(bytes)-1]
	case strings.HasSuffix(bytes, "g"):
		factor = 1024 << 20
		bytes = bytes[0 : len(bytes)-1]
	}
	var val int64
	val, err = strconv.ParseInt(bytes, 10, 64)
	if err != nil {
		return
	}
	return int(val) * factor, err
}

func humanByteRate(bytesPerSec float64) string {
	switch {
	case bytesPerSec < KB:
		return fmt.Sprintf("%0.2f bytes/sec", bytesPerSec)
	case bytesPerSec < MB:
		return fmt.Sprintf("%0.2f KB/sec", bytesPerSec/KB)
	case bytesPerSec < GB:
		return fmt.Sprintf("%0.2f MB/sec", bytesPerSec/MB)
	default:
		return fmt.Sprintf("%0.2f GB/sec", bytesPerSec/GB)
	}
}

func run(file *os.File) (written int, err error) {
	w := bufio.NewWriter(file)
	bs := make([]byte, 4096)
	numBytes, err := parseHumanBytes(*bytes)
	if err != nil {
		return written, err
	}
	for written < numBytes {
		sz := numBytes - written
		if sz > len(bs) {
			sz = len(bs)
		}
		n, err := w.Write(bs[0:sz])
		if err != nil {
			return written, err
		}
		written += n
	}
	err = w.Flush()
	return
}
