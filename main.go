package main

import (
	"flag"
	"fmt"
	"github.com/ginhelper/project"
	"os"
)

func main() {

	currentPath, _ := os.Getwd()
	fmt.Println(currentPath)
	flag.Parse()

	arg := flag.Args()
	if arg[0] == "create" {
		project.GenerateProject()
	}
}
