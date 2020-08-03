package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"os"
	"reflect"
	"time"
	"unsafe"
)

type void struct{}

// bytesToString converts byte to string without copy
func bytesToString(bytes []byte) (s string) {
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	str := (*reflect.StringHeader)(unsafe.Pointer(&s))
	str.Data = slice.Data
	str.Len = slice.Len
	return s
}

// readLine reads a line from a buffered reader
func readLine(reader *bufio.Reader) ([]byte, error) {
	line := []byte{}
	for true {
		_line, isPrefix, err := reader.ReadLine()
		if err != nil {
			return line, err
		}
		line = append(line, _line...)
		if !isPrefix {
			break
		}
	}
	return line, nil
}

func main() {
	hashes := make(map[string]void)
	var member void

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	startTime := time.Now()
	nline := 0
	nunique := 0
	for true {
		// Read line
		line, err := readLine(reader)
		if err != nil {
			break
		}
		nline++

		// Log progress
		if nline%100000 == 0 {
			fmt.Fprintf(os.Stderr, "\rRead %d lines", nline)
		}

		// Hash line
		h := md5.Sum(line)
		hstring := bytesToString(h[:])

		// Check if new line has been seen before
		_, exists := hashes[hstring]
		if exists {
			continue
		}

		nunique++
		hashes[hstring] = member
		writer.Write(line)
		writer.WriteByte('\n')
	}
	writer.Flush()

	timeTaken := time.Now().Sub(startTime)
	sentsPerSecond := float64(nline) / timeTaken.Seconds()
	fmt.Fprintf(os.Stderr, "\rRead %d lines.\n", nline)
	fmt.Fprintf(os.Stderr, "Found %d unique.\n", nunique)
	fmt.Fprintf(os.Stderr, "Done in %v.\n", timeTaken)
	fmt.Fprintf(os.Stderr, "%.2f sents/s\n", sentsPerSecond)
}
