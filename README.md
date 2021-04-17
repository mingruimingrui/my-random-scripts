**my-random-scripts** is a personal collection of random scripts that I use on a
semi-frequent basis.

My work is that of a NLP software engineer. Most of these scripts are used for
general purpose processing/formatting through large amounts of text data.

## `re2`
> CLI for re2 in golang

`grep`, `awk`, `sed` all uses different syntax but I'm most familiar with PCRE
style regex.
The implementation of this script is single threaded.

### Input/Output
This program uses `stdin` and `stdout`.
Progress and verbosity will be sent to `stderr`.

### Usage
```
Usage: re2 [-ifar] [-e] [-v] [-p] pattern [replacement]
  -a	Find all substrings
  -e	Unescape replacement string
  -f	Find substrings
  -i	Ignore lines containing pattern
  -p	Progress?
  -r	Replace all substrings
  -v	Verbose?
```

## `deduplicate-lines`
> Text file deduplication

Internally, this script hashes each line using md5 and keeps that in memory.
The implementation of this script is single threaded.

### Input/Output
This program uses `stdin` and `stdout`.
Progress and verbosity will be sent to `stderr`.

### Usage
```
Usage: ./bin/deduplicate-lines [-v] [-p]
  -p	Progress?
  -v	Verbose?
```

## `deduplicate-csv`
> CSV file deduplication

The concept is similar to deduplicate lines but for CSV but perhaps you only
want to deduplicate on certain columns.
Internally, this script hashes each record using md5 and keeps that in memory.
The implementation of this script is single threaded.

### Input/Output
This program uses `stdin` and `stdout`.
Progress and verbosity will be sent to `stderr`.

### Usage
```
Usage: ./bin/deduplicate-csv [-s SEP] [-v] [-p] col_num ...
  -p	Progress?
  -s string
    	Separator to be used in key (default "|||")
  -v	Verbose?
```
