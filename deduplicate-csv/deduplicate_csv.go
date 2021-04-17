package main

import (
	"crypto/md5"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

type void struct{}

var sepFlag *string = flag.String("s", "|||", "Separator to be used in key")
var verboseFlag *bool = flag.Bool("v", false, "Verbose?")
var progressFlag *bool = flag.Bool("p", false, "Progress?")

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-s SEP] [-v] [-p] col_num ...\n", os.Args[0])
	flag.PrintDefaults()
}

// stringToBytes converts string to bytes without copy
func stringToBytes(s string) (bytes []byte) {
	str := (*reflect.StringHeader)(unsafe.Pointer(&s))
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	slice.Data = str.Data
	slice.Len = str.Len
	return bytes
}

// bytesToString converts bytes to string without copy
func bytesToString(bytes []byte) (s string) {
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	str := (*reflect.StringHeader)(unsafe.Pointer(&s))
	str.Data = slice.Data
	str.Len = slice.Len
	return s
}

func parseIntSlice(arr []string) []int {
	var slice []int
	for _, value := range arr {
		num, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal("Expecting column numbers to be integer.")
		}
		slice = append(slice, num)
	}
	return slice
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("Require atleast 1 column number.")
	}
	colsToKeep := parseIntSlice(flag.Args())
	if *verboseFlag {
		fmt.Println(colsToKeep)
	}
	sep := stringToBytes(*sepFlag)

	hashes := make(map[string]void)
	var member void

	reader := csv.NewReader(os.Stdin)
	writer := csv.NewWriter(os.Stdout)

	startTime := time.Now()
	nrecord := 0
	nunique := 0

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		nrecord++

		if *progressFlag && (nrecord%100000 == 0) {
			fmt.Fprintf(os.Stderr, "\rRead %d records", nrecord)
		}

		var key []byte
		for i, colNum := range colsToKeep {
			value := record[colNum]
			if i > 0 {
				key = append(key, sep...)
			}
			key = append(key, stringToBytes(value)...)
		}

		h := md5.Sum(key)
		hstring := bytesToString(h[:])

		_, exists := hashes[hstring]
		if exists {
			continue
		}

		nunique++
		hashes[hstring] = member
		writer.Write(record)
	}
	writer.Flush()

	if *verboseFlag {
		timeTaken := time.Since(startTime)
		recordsPerSecond := float64(nrecord) / timeTaken.Seconds()
		fmt.Fprintf(os.Stderr, "\rRead %d records.\n", nrecord)
		fmt.Fprintf(os.Stderr, "Found %d unique.\n", nunique)
		fmt.Fprintf(os.Stderr, "Done in %v.\n", timeTaken)
		fmt.Fprintf(os.Stderr, "%.2f records/s\n", recordsPerSecond)
	} else if *progressFlag {
		fmt.Fprint(os.Stderr, "\n")
	}
}
