package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var ignoreFlag *bool = flag.Bool("i", false, "Ignore lines containing pattern")
var findFlag *bool = flag.Bool("f", false, "Find substrings")
var findAllFlag *bool = flag.Bool("a", false, "Find all substrings")
var replaceFlag *bool = flag.Bool("r", false, "Replace all substrings")
var unescapeFlag *bool = flag.Bool("e", false, "Unescape replacement string")
var verboseFlag *bool = flag.Bool("v", false, "Verbose?")
var progressFlag *bool = flag.Bool("p", false, "Progress?")

var Usage = func() {
	fmt.Fprintf(
		os.Stderr,
		"Usage: %s "+
			"[-ifar] [-e] [-v] [-p] "+
			"pattern [replacement]\n",
		os.Args[0],
	)
	flag.PrintDefaults()
}

type methodType string

const (
	searchMethod  methodType = "search"
	ignoreMethod  methodType = "ignore"
	findMethod    methodType = "find"
	findAllMethod methodType = "findAll"
	replaceMethod methodType = "replace"
)

func checkIsTTY() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice != 0
}

func unscapeString(text string) string {
	text = strings.ReplaceAll(text, "\\a", "\a")
	text = strings.ReplaceAll(text, "\\b", "\b")
	text = strings.ReplaceAll(text, "\\f", "\f")
	text = strings.ReplaceAll(text, "\\n", "\n")
	text = strings.ReplaceAll(text, "\\r", "\r")
	text = strings.ReplaceAll(text, "\\t", "\t")
	text = strings.ReplaceAll(text, "\\v", "\v")
	return text
}

// readLine reads a line from a buffered reader
func readLine(reader *bufio.Reader) ([]byte, error) {
	line := []byte{}
	for {
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
	flag.Usage = Usage
	flag.Parse()

	// Determine method, patternStr, replaceStr
	if flag.NArg() < 1 {
		log.Fatal("Pattern not provided")
	}
	patternStr := flag.Arg(0)
	replaceStr := ""

	var method methodType
	if *ignoreFlag {
		method = ignoreMethod
	} else if *findFlag {
		method = findMethod
	} else if *findAllFlag {
		method = findAllMethod
	} else if *replaceFlag {
		method = replaceMethod
		if flag.NArg() > 1 {
			replaceStr = flag.Arg(1)
		}
	} else {
		method = searchMethod
	}

	// Compile pattern
	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		log.Fatal(err)
	}

	// Format replaceStr
	if *unescapeFlag {
		replaceStr = unscapeString(replaceStr)
	}
	replaceBytes := []byte(replaceStr)

	if *verboseFlag {
		fmt.Fprintf(os.Stderr, "Method: %v\n", method)
		fmt.Fprintf(os.Stderr, "Pattern: %v\n", pattern)
		if method == replaceMethod {
			fmt.Fprintf(os.Stderr, "Replace str: %v\n", replaceStr)
		}
	}

	isTTY := checkIsTTY()

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	startTime := time.Now()
	nline := 0

	writeLine := func(text []byte) {
		writer.Write(text)
		writer.WriteByte('\n')
		if isTTY {
			writer.Flush()
		}
	}

	for {
		line, err := readLine(reader)
		if err != nil {
			break
		}
		nline++

		if *progressFlag && (nline%100000 == 0) {
			fmt.Fprintf(os.Stderr, "\rRead %d lines", nline)
		}

		switch method {

		case searchMethod:
			if pattern.Match(line) {
				writeLine(line)
			}

		case ignoreMethod:
			if !pattern.Match(line) {
				writeLine(line)
			}

		case findMethod:
			foundStr := pattern.Find(line)
			if len(foundStr) > 0 {
				writeLine(foundStr)
			}

		case findAllMethod:
			for _, foundStr := range pattern.FindAll(line, -1) {
				if len(foundStr) > 0 {
					writeLine(foundStr)
				}
			}

		case replaceMethod:
			line = pattern.ReplaceAll(line, replaceBytes)
			writeLine(line)

		}

	}
	writer.Flush()

	if *verboseFlag {
		timeTaken := time.Since(startTime)
		sentsPerSecond := float64(nline) / timeTaken.Seconds()
		fmt.Fprintf(os.Stderr, "\rRead %d lines.\n", nline)
		fmt.Fprintf(os.Stderr, "Done in %v.\n", timeTaken)
		fmt.Fprintf(os.Stderr, "%.2f sents/s\n", sentsPerSecond)
	} else if *progressFlag {
		fmt.Fprint(os.Stderr, "\n")
	}
}
