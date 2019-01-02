package gweb

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Args reads stanard input and converts each line to an array of strings.
// TODO: safe split with quotes, etc.
func Args() <-chan []string {
	chargs := make(chan []string)
	go func() {
		sc := bufio.NewScanner(Stdin)
		for sc.Scan() {
			var args []string
			line := strings.TrimSpace(sc.Text())
			if line != "" {
				args = strings.Split(line, " ")
			}
			chargs <- args
		}
		if err := sc.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}()
	return chargs
}
