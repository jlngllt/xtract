package main

import (
	"bufio"
	"fmt"
	"os"

	"mvdan.cc/xurls/v2"
)

func main() {
	rx := xurls.Relaxed()
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		text := sc.Text()
		urlLoc := rx.FindAllIndex([]byte(text), -1)
		for _, loc := range urlLoc {
			// check for wildcard,... not so good code
			if loc[0] > 2 && text[loc[0]-2:loc[0]] == "*." {
				fmt.Println(text[loc[0]-2 : loc[1]])
			} else {
				fmt.Println(text[loc[0]:loc[1]])
			}
		}
	}
}
