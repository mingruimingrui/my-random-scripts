**my-random-scripts** is a personal collection of random scripts that I use on a
semi-frequent basis.

My work is that of a NLP software engineer. Most of these scripts are used for
general purpose processing/formatting through large amounts of text data.

## `re2`

`grep`, `awk`, `sed` all uses different syntax but I'm most familiar with PCRE style
regex. This script is basiscally re2 implemented in golang but bought to the console.

### Options

| Flag | Description |
| ---- | ----------- |
| `-i` | Ignore lines containing pattern |
| `-f` | Find substrings |
| `-a` | Find all substrings |
| `-r` | Replace all substrings |
| `-e` | Unescape replacement string (else `\n` will be interpreted literally) |
| `-v` | Verbose? |
| `-p` | Progress? |

### Input/Output

This program uses `stdin` and `stdout`.
Progress and verbosity will be sent to `stderr`.
