package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Options []string

func (v Options) String() string {
	return "[" + strings.Join(v, ", ") + "]"
}

func (v Options) Detailed() string {
	var res string
	for _, s := range v {
		res += fmt.Sprintf("- %s\n", s)
	}
	res += "- exit\n"

	return res
}

func readInputAsString(req ...string) string {
	for _, v := range req {
		fmt.Print(v)
	}
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input.Text()
}

func readInputAsInt(req ...string) int {
	s := readInputAsString(req...)
	v, _ := strconv.Atoi(s)
	return v
}
func consoleHelper(subOps Options, handlers ...func() bool) bool {
	for {
		fmt.Print(subOps.Detailed())
		fmt.Print(": ")
		s := readInputAsString()
		if s == "exit" {
			// clear stdout
			fmt.Print("\033[H\033[2J")
			return true
		}

		matchI := -1
		for i, sop := range subOps {
			if s == sop {
				matchI = i
				break
			}
		}

		if matchI < 0 {
			fmt.Printf("operation not found. available operation %s\n", subOps)
		} else if len(handlers) > matchI {
			// clear stdout
			fmt.Print("\033[H\033[2J")

			for {
				// print 'breadcrumbs'
				fmt.Printf("# %s\n", subOps[matchI])

				// run handler
				if handlers[matchI]() {
					break
				}
			}
			break
		} else {
			fmt.Printf("handler not found for %s\n", subOps[matchI])
			break
		}
	}
	return false
}
