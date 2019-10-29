package main

import (
	"flag"
	"fmt"
	"github.com/will110/ginhelper/project"
)

const usage = `
Ginhelper is a gin web framework helper, Help you develop fast.

USAGE

    ginhelper [arguments]

AVAILABLE ARGUMENTS

	create -> Create single gin web framework
`

func main() {
	flag.Parse()
	arg := flag.Args()
	if len(arg) == 0 {
		fmt.Println(usage)
		return
	}

	if arg[0] == "create" {
		project.GenerateProject()
	} else {
		fmt.Println(usage)
	}
}
