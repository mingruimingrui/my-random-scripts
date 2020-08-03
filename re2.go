package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"time"
	"unsafe"
)

type methodType string

const (
	searchMethod  methodType = "search"
	ignoreMethod  methodType = "ignore"
	findMethod    methodType = "find"
	findAllMethod methodType = "findAll"
	replaceMethod methodType = "replace"
)

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

// assertEqual checks if 2 values are the same
func assertEqual(a interface{}, b interface{}, msg string) {
	if a == b {
		return
	}
	if len(msg) == 0 {
		msg = fmt.Sprintf("%v != %v", a, b)
	}
	log.Fatal(msg)
}

func main() {
	var patternStr string
	var replaceStr []byte

	// Flags
	var method methodType
	ignoreFlag := flag.Bool("i", false, "Ignore lines with pattern")
	findFlag := flag.Bool("f", false, "Search and find substring")
	findAllFlag := flag.Bool("a", false, "Search and find all substrings")
	replaceFlag := flag.Bool("r", false, "Replace pattern")
	verboseFlag := flag.Bool("v", false, "Verbose")
	flag.Parse()

	// Determine patternStr and replaceStr based on flag
	if *ignoreFlag {
		assertEqual(flag.NArg(), 1, "Expecting (only) 1 positional argument")
		method = "ignore"
	} else if *findFlag {
		assertEqual(flag.NArg(), 1, "Expecting (only) 1 positional argument")
		method = "find"

	} else if *findAllFlag {
		assertEqual(flag.NArg(), 1, "Expecting (only) 1 positional argument")
		method = "findAll"

	} else if *replaceFlag {
		assertEqual(flag.NArg(), 2, "Expecting (only) 2 positional argument")
		method = "replace"
		replaceStr = []byte(flag.Arg(1))

	} else {
		// Default case is search like in grep
		assertEqual(flag.NArg(), 1, "Expecting (only) 1 positional argument")
		method = "search"
	}
	patternStr = flag.Arg(0)

	// Compile pattern
	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		log.Fatal(err)
	}

	if *verboseFlag {
		fmt.Fprintf(os.Stderr, "Method: %v\n", method)
		fmt.Fprintf(os.Stderr, "Pattern: %v\n", pattern)
		if method == replaceMethod {
			fmt.Fprintf(os.Stderr, "Replace str: %v\n", replaceStr)
		}
	}

	// Ensure is piped
	fi, _ := os.Stdin.Stat()
	if fi.Mode()&os.ModeCharDevice != 0 {
		log.Fatal("Data must be piped")
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	startTime := time.Now()
	nline := 0

	for true {
		line, err := readLine(reader)
		if err != nil {
			break
		}
		nline++

		// Log progress
		if nline%100000 == 0 {
			fmt.Fprintf(os.Stderr, "\rRead %d lines", nline)
		}

		switch method {
		case ignoreMethod:
			if !pattern.Match(line) {
				writer.Write(line)
				writer.WriteByte('\n')
			}

		case searchMethod:
			if pattern.Match(line) {
				writer.Write(line)
				writer.WriteByte('\n')
			}

		case findMethod:
			foundStr := pattern.Find(line)
			if len(foundStr) > 0 {
				writer.Write(foundStr)
				writer.WriteByte('\n')
			}

		case findAllMethod:
			for _, foundStr := range pattern.FindAll(line, -1) {
				writer.Write(foundStr)
				writer.WriteByte('\n')
			}

		case replaceMethod:
			line = pattern.ReplaceAll(line, replaceStr)
			writer.Write(line)
			writer.WriteByte('\n')

		}

	}
	writer.Flush()

	timeTaken := time.Now().Sub(startTime)
	sentsPerSecond := float64(nline) / timeTaken.Seconds()
	fmt.Fprintf(os.Stderr, "\rRead %d lines.\n", nline)
	fmt.Fprintf(os.Stderr, "Done in %v.\n", timeTaken)
	fmt.Fprintf(os.Stderr, "%.2f sents/s\n", sentsPerSecond)
}